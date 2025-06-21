package client

import (
	"context"
	"fmt"
	"path"

	gcConfig "github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
)

// SQSClient wraps the GoCommons SQS Publisher for CreateBulkOrder
type SQSClient struct {
	Publisher *sqs.Publisher
}

// NewSQSClient initializes the SQS publisher using config
func NewSQSClient(ctx context.Context) (*SQSClient, error) {
	// Read queue URL from config
	queueURL := gcConfig.GetString(ctx, "sqs.bulk_order_queue_url")
	if queueURL == "" {
		return nil, fmt.Errorf("NewSQSClient: missing sqs.bulk_order_queue_url in config")
	}

	queueName := path.Base(queueURL)

	// Build SQS config from config.yaml
	sqsCfg := &sqs.Config{
		Account:  gcConfig.GetString(ctx, "sqs.account"),
		Endpoint: gcConfig.GetString(ctx, "sqs.endpoint"),
		Region:   gcConfig.GetString(ctx, "sqs.region"),
	}

	log.Infof("👉 SQS config: URL=%s endpoint=%s region=%s account=%s", queueURL, sqsCfg.Endpoint, sqsCfg.Region, sqsCfg.Account)

	// Initialize queue
	queue, err := sqs.NewStandardQueue(ctx, queueName, sqsCfg)
	if err != nil {
		log.DefaultLogger().Errorf("NewSQSClient: failed to create SQS queue: %v", err)
		return nil, err
	}

	// Create publisher
	publisher := sqs.NewPublisher(queue)

	log.DefaultLogger().Infof("✅ SQS Publisher initialized for queue: %s", queueName)

	return &SQSClient{
		Publisher: publisher,
	}, nil
}

// PublishCreateBulkOrderEvent sends a message payload to the SQS queue
func (c *SQSClient) PublishCreateBulkOrderEvent(ctx context.Context, payload []byte) error {
	msg := &sqs.Message{
		Value: payload,
	}

	if err := c.Publisher.Publish(ctx, msg); err != nil {
		log.DefaultLogger().Errorf("❌ SQS publish failed: %v", err)
		return err
	}

	log.DefaultLogger().Infof("✅ SQS message published successfully")
	return nil
}