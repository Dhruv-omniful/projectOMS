package main

import (
	"fmt"
	"time"

	"github.com/omniful/go_commons/config"
	// "github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/health"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"

	"ims/postgres"
	"ims/router"
)

func main() {
	// Initialize configuration with hot-reload every 30 seconds
	if err := config.Init(30 * time.Second); err != nil {
		log.Panicf("Failed to initialize config: %v", err)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Failed to create config context: %v", err)
	}

	// Initialize Postgres and run migrations
	pr.InitPostgres(ctx)

	// Initialize Redis (if Redis code is similar)
	pr.InitRedis(ctx)

	// Set log level from config
	log.SetLevel(config.GetString(ctx, "log.level"))

	port := config.GetInt(ctx, "server.port")
	log.Infof("Starting IMS server on port %d", port)

	level := config.GetString(ctx, "log.level")
	log.SetLevel(level)

	logOpts := http.LoggingMiddlewareOptions{
		Format:      config.GetString(ctx, "log.format"),
		Level:       level,
		LogRequest:  true,
		LogResponse: true,
		LogHeader:   false,
	}

	// Initialize HTTP server with common middleware
	server := http.InitializeServer(
		fmt.Sprintf(":%d", port),
		config.GetDuration(ctx, "server.readTimeout"),
		config.GetDuration(ctx, "server.writeTimeout"),
		config.GetDuration(ctx, "server.idleTimeout"),
		false, // TLS disabled
		env.RequestID(),
		env.Middleware(config.GetString(ctx, "env")),
		config.Middleware(),
		http.RequestLogMiddleware(logOpts),
	)

	// Health check endpoint
	server.Engine.GET("/health", health.HealthcheckHandler())

	// Register application routes
	routes.RegisterRoutes(server.Engine)

	// Start server (blocking)
	if err := server.StartServer("IMS"); err != nil {
		log.Errorf("Server shutdown with error: %v", err)
	}
}
