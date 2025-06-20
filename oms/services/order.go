package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/log"

	"github.com/dhruv/oms/client"
	"github.com/dhruv/oms/model"
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
	log.Infof("🔍 Validating S3 path: %s", s3Path)

	_, err := s.S3.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.S3.Bucket,
		Key:    &s3Path,
	})
	if err != nil {
		log.Errorf("❌ S3 HeadObject failed: %v", err)
		return fmt.Errorf("failed to validate S3 path %s: %w", s3Path, err)
	}

	log.Infof("✅ S3 file exists: %s", s3Path)

	// Build CreateBulkOrderEvent
	event := model.CreateBulkOrderEvent{
		S3Path:     s3Path,
		TenantID:   "tenant-123", // Replace with real tenant ID if available
		UploadedAt: time.Now().Format(time.RFC3339),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Errorf("❌ Failed to marshal CreateBulkOrderEvent: %v", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish event to SQS
	if err := s.SQSClient.PublishCreateBulkOrderEvent(ctx, payload); err != nil {
		log.Errorf("❌ Failed to publish CreateBulkOrderEvent: %v", err)
		return fmt.Errorf("failed to publish SQS event: %w", err)
	}

	log.Infof("✅ CreateBulkOrderEvent published to SQS")

	return nil
}
