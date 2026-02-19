package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"quantatomai/grid-service/compute"
	"quantatomai/grid-service/fetcher"
	"quantatomai/grid-service/handlers"
	"quantatomai/grid-service/mapping"
	"quantatomai/grid-service/planner"
	"quantatomai/grid-service/storage"
)

func main() {
	// 1. Core Infrastructure Initialization
	db := initDB()
	defer db.Close()

	rdb := initRedis()
	defer rdb.Close()

	// 2. Metadata & Planning Layer
	// We use the cached decorator to avoid slamming Postgres for repeated axis lookups.
	pgMeta, err := mapping.NewPostgresMetadataResolver(db, "default_model", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to init postgres metadata: %v", err)
	}
	defer pgMeta.Close()

	cachedMeta := mapping.NewCachedMetadataResolver(pgMeta, 10*time.Minute)
	planr := planner.NewPlanner(cachedMeta)

	// 3. Fetching Layer (Hot Tier)
	// Resilient Redis fetcher with concurrent workers and circuit breaker.
	atomFetcher := fetcher.NewRedisAtomFetcher(rdb, fetcher.RedisFetcherConfig{
		ChunkSize:        1000,
		WorkerCount:      10,
		FailureThreshold: 5,
		OpenDuration:     10 * time.Second,
	})

	// 4. Compute Layer (FX & Multi-Axial)
	// We use the postgres resolver as the currency provider (implements CurrencyMetadataProvider).
	// In production, we'd use a dedicated rate provider (e.g. from a market data service).
	currResolver := compute.NewCurrencyResolverMetadata(pgMeta, []compute.CurrencyDimensionBinding{
		{Role: compute.DimensionRoleEntity, DimensionID: 101, Index: 0}, // Row 0 = Entity
	})

	// Placeholder FX provider (would hit fx_rates table)
	fxProvider := &stubFXProvider{} 

	fxTransformer := compute.NewFXTransformer(compute.FXConfig{
		TargetCurrency: "USD",
		EnableAudit:    true,
	}, fxProvider, currResolver)

	compEngine := compute.NewDefaultComputeEngine([]compute.AtomTransformer{
		fxTransformer,
	})

	// 5. Grid Storage Subsystem (Ultra-Diamond)
	cfg := storage.NewDefaultCacheConfig()
	cfg.WireFormat = storage.WireFormatFlatBuffer
	cfg.EnableChunkedHydration = true

	// nodeID for orchestration; in reality, use os.Hostname()
	nodeID := "grid-node-master"

	tieredCache, _, err := storage.NewGridSubsystem(
		context.Background(),
		rdb,
		nodeID,
		cfg,
		nil, // prefetchRecomputeFn
		nil, // coordinator
		nil, // xfetchRefresher
	)
	if err != nil {
		log.Fatalf("failed to init grid storage: %v", err)
	}

	// 6. Handler Orchestration
	gridHandler := handlers.NewGridQueryHandler(
		planr,
		atomFetcher,
		compEngine,
		nil,         // CircuitBreaker (stubbed)
		tieredCache,
	)

	// 7. HTTP Routing
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Core API
	router.POST("/grid/query", gridHandler.HandleGridQuery)
	router.POST("/grid/writeback", handlers.HandleWriteback)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("QuantatomAI Grid Service starting on :%s (WireFormat: %s)", port, cfg.WireFormat)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run router: %v", err)
	}
}

func initDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/quantatomai?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return db
}

func initRedis() *redis.Client {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		addr = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

// stubFXProvider for initialization demo
type stubFXProvider struct{}

func (s *stubFXProvider) GetRate(ctx context.Context, src, tgt string, asOf time.Time) (float64, error) {
	if src == tgt {
		return 1.0, nil
	}
	return 1.1, nil // Static demo rate
}
