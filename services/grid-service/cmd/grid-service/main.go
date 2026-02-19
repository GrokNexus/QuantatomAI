// File: cmd/grid-service/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"quantatomai/grid-service/src/handlers"
	"quantatomai/grid-service/src/projection"
	"quantatomai/grid-service/src/storage"
)

type RedisClientImpl struct{}

func (r *RedisClientImpl) Set(ctx context.Context, key string, value []byte) error {
	return nil
}

func (r *RedisClientImpl) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func main() {
	storage.GlobalArenaManager = projection.NewArenaManager(64 * 1024 * 1024)

	redisClient := &RedisClientImpl{}
	gridCache := storage.NewRedisGridCache(redisClient)
	gridHandler := handlers.NewGridQueryHandler(gridCache)

	mux := http.NewServeMux()
	mux.HandleFunc("/grid/query", gridHandler.HandleGridQuery)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Grid Service running on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
