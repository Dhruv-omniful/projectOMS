package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"

	"github.com/dhruv/oms/client"
)

// OrderService handles order-related logic
type OrderService struct {
	S3        *client.S3Client
	SQSClient *client.SQSClient
}

// NewOrderService creates a new OrderService
func NewOrderService(s3Client *client.S3Client, sqsClient *client.SQSClient) *OrderService {
	return &OrderService{
		S3:        s3Client,
		SQSClient: sqsClient,
	}
}

// ProcessCSV validates S3 path and pushes SQS event
func (s *OrderService) ProcessCSV(ctx context.Context, s3Path string) error {
	log.Infof(" Validating S3 path: %s", s3Path)

	_, err := s.S3.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.S3.Bucket,
		Key:    &s3Path,
	})
	if err != nil {
		log.Errorf(" S3 HeadObject failed: %v", err)
		return fmt.Errorf("failed to validate S3 path %s: %w", s3Path, err)
	}

	log.Infof(" S3 file exists: %s", s3Path)

	//  Create payload that CSV worker expects
	payload := map[string]string{
		"Bucket": config.GetString(ctx, "s3.bucket"),
		"Key":    s3Path,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorf(" Failed to marshal SQS payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if err := s.SQSClient.PublishCreateBulkOrderEvent(ctx, data); err != nil {
		log.Errorf(" Failed to publish SQS event: %v", err)
		return fmt.Errorf("failed to publish event to SQS: %w", err)
	}

	log.Infof(" CreateBulkOrderEvent published to SQS: %v", payload)
	return nil
}
