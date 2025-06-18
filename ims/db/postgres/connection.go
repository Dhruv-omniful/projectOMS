package db

import (
	"log"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/postgres"

	"ims/context" // ✅ This should match your actual context package name
)

func ConnectDB() *postgres.DbCluster {
	ctx := imscontext.GetContext()

	masterConfig := postgres.DBConfig{
		Host:                   config.GetString(ctx, "postgres.master.host"),
		Port:                   config.GetString(ctx, "postgres.master.port"),
		Username:               config.GetString(ctx, "postgres.master.username"),
		Password:               config.GetString(ctx, "postgres.master.password"),
		Dbname:                 config.GetString(ctx, "postgres.master.dbname"),
		MaxOpenConnections:      config.GetInt(ctx, "postgres.master.max_open_conns"),
		MaxIdleConnections:      config.GetInt(ctx, "postgres.master.max_idle_conns"),
		ConnMaxLifetime:         config.GetDuration(ctx, "postgres.master.conn_max_lifetime"),
		DebugMode:               config.GetBool(ctx, "postgres.master.debug_mode"),
		PrepareStmt:             true,
		SkipDefaultTransaction:  true,
	}

	slavesConfig := make([]postgres.DBConfig, 0)

	log.Println("Initializing DB connection...")

	dbCluster := postgres.InitializeDBInstance(masterConfig, &slavesConfig)

	log.Println("✅ Database connection established")

	return dbCluster
}
