package redis

import (
	"time"

	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
)

var Client *redis.Client

// Start initializes the Redis client (singleton pattern)
func Start() {
	if Client != nil {
		log.Warn("⚠️ Redis client already initialized, skipping re-init")
		return
	}

	config := &redis.Config{
		Hosts:        []string{"localhost:6379"},
		PoolSize:     50,
		MinIdleConn:  10,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  600 * time.Second,
	}

	Client = redis.NewClient(config)
	log.Infof("✅ Redis client initialized")
}

// GetClient safely returns the initialized client
func GetClient() *redis.Client {
	if Client == nil {
		log.Panic("❌ Redis client is not initialized. Call redis.Start() first.")
	}
	return Client
}
