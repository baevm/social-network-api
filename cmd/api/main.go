package main

import (
	"social-network-api/cfg"
	"social-network-api/internal/db"
	"social-network-api/internal/http"
	"social-network-api/internal/redis"
	"social-network-api/pkg/logger"
)

// @title           Social network API
// @version         1.0
// @description     Twitter like api made with golang.

// @host      localhost:5000
// @BasePath  /v1
func main() {
	logger := logger.New()
	err := cfg.Load(".")

	if err != nil {
		logger.Fatalf("Error reading config: %s", err)
	}

	db, err := db.New()

	if err != nil {
		logger.Fatalf("Error starting db: %s", err)
	}

	defer db.Close()

	cache := redis.New(cfg.Get().Redis.Host, cfg.Get().Redis.Port, cfg.Get().Redis.Pass)

	defer cache.Close()

	httpServer := http.New(logger, db, cache)
	httpServer.Run()
}
