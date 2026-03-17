package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	migrate "quantatomai/grid-service/sql/schema"
	"quantatomai/grid-service/src/compute"
	"quantatomai/grid-service/src/fetcher"
	"quantatomai/grid-service/src/handlers"
	"quantatomai/grid-service/src/mapping"
	"quantatomai/grid-service/src/planner"
)

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "run migrations and exit")
	flag.Parse()

	// 1. Core Infrastructure Initialization
	db := initDB()
	defer func() { _ = db.Close() }()

	rdb := initRedis()
	defer func() { _ = rdb.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := migrate.Run(ctx, db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// 2. Metadata & Planning Layer
	// Model ID is configurable to support multiple apps/models per environment.
	modelID := os.Getenv("MODEL_ID")
	if modelID == "" {
		modelID = "default_model"
	}

	// We use the cached decorator to avoid slamming Postgres for repeated axis lookups.
	pgMeta, err := mapping.NewPostgresMetadataResolver(db, modelID, 2*time.Second)
	if err != nil {
		log.Fatalf("failed to init postgres metadata: %v", err)
	}
	defer func() { _ = pgMeta.Close() }()

	if *migrateOnly {
		log.Println("migrations applied; exiting per migrate-only flag")
		return
	}

	cachedMeta := mapping.NewCachedMetadataResolver(pgMeta, 10*time.Minute, 1000)
	planr := planner.NewPlanner(cachedMeta)

	// 3. Fetching Layer (Hot Tier)
	// Resilient Redis fetcher with concurrent workers and circuit breaker.
	atomFetcher := fetcher.NewRedisAtomFetcher(rdb, "atom-cache:", 5*time.Second)

	// 4. Compute Layer (FX & Multi-Axial)
	// We use the postgres resolver as the currency provider (implements CurrencyMetadataProvider).
	// In production, we'd use a dedicated rate provider (e.g. from a market data service).
	currResolver := compute.NewCurrencyResolverMetadata(pgMeta, compute.CurrencyResolverConfig{
		Bindings: []compute.CurrencyDimensionBinding{
			{Role: compute.DimensionRoleEntity, DimensionID: 101, Index: 0}, // Row 0 = Entity
		},
	})

	// Placeholder FX provider (would hit fx_rates table)
	fxProvider := &stubFXProvider{}

	fxTransformer := compute.NewFXTransformer(compute.FXConfig{
		TargetCurrency: "USD",
		EnableAudit:    true,
	}, fxProvider, currResolver)

	compEngine := compute.NewDefaultComputeEngine(
		compute.WithTransformers(fxTransformer),
	)

	// 5. Grid Storage Subsystem (Deleted/Legacy Tiered Cache)

	// 6. Handler Orchestration
	gridHandler := handlers.NewGridHandler(
		planr,
		atomFetcher,
		compEngine,
		currResolver,
	)
	metadataHandler := handlers.NewMetadataHandler(planr.MetadataResolver())
	metadataGraphHandler := handlers.NewMetadataGraphHandler(db)

	// 7. HTTP Routing
	router := gin.Default()

	// CORS for UI (configure via CORS_ORIGINS, comma-separated; default allows localhost:3000)
	router.Use(func(c *gin.Context) {
		origins := os.Getenv("CORS_ORIGINS")
		if origins == "" {
			origins = "http://localhost:3000"
		}
		allowed := strings.Split(origins, ",")
		reqOrigin := c.GetHeader("Origin")
		for _, o := range allowed {
			if strings.TrimSpace(o) == reqOrigin {
				c.Header("Access-Control-Allow-Origin", reqOrigin)
				c.Header("Vary", "Origin")
				break
			}
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	})

	// 7.1 Git-Flow Metadata Routing (Layer 8.1)
	branchHandler := handlers.NewBranchHandler(db)
	branchHandler.RegisterRoutes(router)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Metadata discovery
	metadataHandler.RegisterRoutes(router)
	metadataGraphHandler.RegisterRoutes(router)

	// Core API
	router.POST("/grid/query", gridHandler.HandleGridQuery)
	router.POST("/grid/writeback", handlers.HandleWriteback)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("QuantatomAI Grid Service starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run router: %v", err)
	}
}

func initDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://quantatomai:quantatomai@localhost:5432/quantatomai?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
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
