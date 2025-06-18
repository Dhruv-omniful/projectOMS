package imscontext

import (
	"context"
	"log"
	"time"

	"github.com/omniful/go_commons/config"
)

var ctx context.Context

func init() {
	// Mandatory: Initialize config before context
	err := config.Init(10 * time.Second) // Loads YAML config or env
	if err != nil {
		log.Panicf("Error initializing config: %v", err)
	}

	ctx, err = config.TODOContext()
	if err != nil {
		log.Panicf("Error creating context: %v", err)
	}
}

// GetContext returns app context loaded with config
func GetContext() context.Context {
	return ctx
}
