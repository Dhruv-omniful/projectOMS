package worker

import (
	"context"
	"path"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	"github.com/dhruv/oms/client"
)

func StartCSVProcessor(ctx context.Context, imsClient *client.IMSClient) {
	logger := log.DefaultLogger()

	queueURL := config.GetString(ctx, "sqs.bulk_order_queue_url")
	logger.Infof(" Listening for SQS messages on %s", queueURL)

	qName := path.Base(queueURL)

	sqsCfg := &sqs.Config{
		Account:  config.GetString(ctx, "sqs.account"),
		Endpoint: config.GetString(ctx, "sqs.endpoint"),
		Region:   config.GetString(ctx, "sqs.region"),
	}

	qObj, err := sqs.NewStandardQueue(ctx, qName, sqsCfg)
	if err != nil {
		logger.Panicf(" Failed to create SQS queue: %v", err)
	}

	handler := NewQueueHandler(imsClient)

	consumer, err := sqs.NewConsumer(
		qObj,
		uint64(config.GetInt(ctx, "sqs.consumer.worker_count")),
		uint64(config.GetInt(ctx, "sqs.consumer.concurrency_per_worker")),
		handler,
		int64(config.GetInt(ctx, "sqs.consumer.batch_size")),
		int64(config.GetInt(ctx, "sqs.consumer.visibility_timeout")),
		false,
		false,
	)
	if err != nil {
		logger.Panicf(" Failed to start SQS consumer: %v", err)
	}

	consumer.Start(ctx)
}
