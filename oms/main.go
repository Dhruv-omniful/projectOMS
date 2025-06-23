package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/health"
	// commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	// "github.com/omniful/go_commons/redis"

	"github.com/dhruv/oms/api"
	"github.com/dhruv/oms/client"
	"github.com/dhruv/oms/services"
	"github.com/dhruv/oms/worker"
)

func main() {
	// === CONFIG INIT ===
	if err := config.Init(30 * time.Second); err != nil {
		panic(err)
	}

	// === CONTEXT ===
	ctx, err := config.TODOContext()
	if err != nil {
		panic(err)
	}

	// === LOGGER ===
	lvl := config.GetString(ctx, "log.level")
	log.SetLevel(lvl)
	log.Infof("🚀 Starting OMS on port %d", config.GetInt(ctx, "server.port"))

	// === S3 CLIENT ===
	s3Client, err := client.NewS3Client(ctx)
	if err != nil {
		log.Panicf("❌ Failed to initialize S3 client: %v", err)
	}
	log.Info("✅ S3 client initialized successfully")

	fmt.Println("AWS_ACCESS_KEY_ID:", os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Println("AWS_SECRET_ACCESS_KEY:", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	fmt.Println("LOCAL_SQS_ENDPOINT:", os.Getenv("LOCAL_SQS_ENDPOINT"))

	// === SQS CLIENT ===
	sqsClient, err := client.NewSQSClient(ctx)
	if err != nil {
		log.Panicf("❌ Failed to initialize SQS client: %v", err)
	}
	log.Info("✅ SQS client initialized successfully")

	// === KAFKA PRODUCER ===
	client.InitKafkaProducer(ctx)
	log.Info("✅ Kafka producer initialized successfully")

	// === IMS CLIENT SETUP ===

	// imsHTTP, err := commonsHttp.NewHTTPClient(
	// 	"oms-ims-client",
	// 	config.GetString(ctx, "ims.base_url"),
	// 	nil,
	// 	commonsHttp.WithTimeout(5*time.Second),
	// )
	// if err != nil {
	// 	log.Panicf("❌ Failed to init IMS HTTP client: %v", err)
	// }

	imsClient := client.NewIMSClient()

	log.Info("✅ IMS client initialized successfully")

	// === KAFKA CONSUMER (Order Finalizer Worker) ===
	go worker.StartOrderFinalizer(ctx)

	// === ORDER SERVICE ===
	orderService := service.NewOrderService(s3Client, sqsClient)

	// === HANDLERS ===
	handlers := api.NewHandlers(orderService)

	// === START WORKER ===
	go worker.StartCSVProcessor(ctx, imsClient)

	// === SERVER SETUP ===
	port := ":" + strconv.Itoa(config.GetInt(ctx, "server.port"))
	srv := http.InitializeServer(
		port,
		config.GetDuration(ctx, "server.read_timeout"),
		config.GetDuration(ctx, "server.write_timeout"),
		config.GetDuration(ctx, "server.idle_timeout"),
		false,
		env.RequestID(),
		env.Middleware(config.GetString(ctx, "env")),
	)

	// === ROUTES ===
	srv.Engine.GET("/health", health.HealthcheckHandler())
	api.RegisterRoutes(srv.Engine, handlers)

	// === START SERVER ===
	if err := srv.StartServer("oms-service"); err != nil {
		log.Errorf("OMS shutdown error: %v", err)
	}
}
