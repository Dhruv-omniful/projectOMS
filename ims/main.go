package main

import (
	"time"
	// "fmt"
	"log"

	"github.com/omniful/go_commons/http"
	// "github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/migration"

	"ims/db/postgres"
	"ims/context"
	"ims/router"
)

func main() {
	log.Println("🚀 IMS starting...")

	// Recover from panic if any
	defer func() {
		if r := recover(); r != nil {
			log.Println("💥 IMS crashed: %v", r)
		}
	}()

	// Initialize DB
	db := db.ConnectDB()
	if db == nil {
		log.Panic("❌ Database connection failed")
	}
	log.Println("✅ Database connected")

	// Get context
	ctx := imscontext.GetContext()

	// Build DB URL for migrations
	dbURL := migration.BuildSQLDBURL(
		config.GetString(ctx, "postgres.master.host"),
		config.GetString(ctx, "postgres.master.port"),
		config.GetString(ctx, "postgres.master.dbname"),
		config.GetString(ctx, "postgres.master.username"),
		config.GetString(ctx, "postgres.master.password"),
	)
	log.Println("DB URL: %s", dbURL)

	// Path to migrations
	migrationPath := "file://C:/Users/dhruv/Desktop/omni_project/omni_project/ims/migrations"

	// Initialize and apply migrations
	migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
	if err != nil {
		log.Panicf("❌ Migration init failed: %v", err)
	}

	err = migrator.Up()
	if err != nil {
		log.Panicf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migrations applied successfully")

	// Initialize server
	server := http.InitializeServer(
		":8080",
		10*time.Second,  // read timeout
		10*time.Second,  // write timeout
		70*time.Second,  // idle timeout
		false,           // HTTPS (false for now)
		// middlewares if any
	)

	router.Initialize(ctx, server)


	// Start server
	if err := server.StartServer("ims"); err != nil {
		log.Panicf("❌ Server failed to start: %v", err)
	} else {
		log.Println("✅ IMS server is running on :8080")
	}
}
