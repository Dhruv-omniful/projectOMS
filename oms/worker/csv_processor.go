package worker

import (
	"context"
	"path"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
)

// Starts the CSV SQS consumer
func StartCSVProcessor(ctx context.Context) {
	logger := log.DefaultLogger()

	queueURL := config.GetString(ctx, "sqs.bulk_order_queue_url")
	logger.Infof("📥 Listening for SQS messages on %s", queueURL)

	qName := path.Base(queueURL)

	sqsCfg := &sqs.Config{
		Account:  config.GetString(ctx, "sqs.account"),
		Endpoint: config.GetString(ctx, "sqs.endpoint"),
		Region:   config.GetString(ctx, "sqs.region"),
	}

	qObj, err := sqs.NewStandardQueue(ctx, qName, sqsCfg)
	if err != nil {
		logger.Panicf("❌ Failed to create SQS queue: %v", err)
	}

	consumer, err := sqs.NewConsumer(
		qObj,
		uint64(config.GetInt(ctx, "sqs.consumer.worker_count")),
		uint64(config.GetInt(ctx, "sqs.consumer.concurrency_per_worker")),
		&queueHandler{},
		int64(config.GetInt(ctx, "sqs.consumer.batch_size")),
		int64(config.GetInt(ctx, "sqs.consumer.visibility_timeout")),
		false, // autoExtend
		false, // autoDelete
	)
	if err != nil {
		logger.Panicf("❌ Failed to start SQS consumer: %v", err)
	}

	consumer.Start(ctx)
}
