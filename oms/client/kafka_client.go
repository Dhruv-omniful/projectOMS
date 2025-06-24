package client

import (
	"context"
	// "encoding/json"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"

	"github.com/dhruv/oms/model"
)

var kafkaLogger = log.DefaultLogger()

var producer *kafka.ProducerClient

func InitKafkaProducer(ctx context.Context) {
	brokers := config.GetStringSlice(ctx, "kafka.brokers")
	version := config.GetString(ctx, "kafka.version")

	if version == "" {
		log.Panicf(" Kafka version is missing in config")
	}

	producer = kafka.NewProducer(
		kafka.WithBrokers(brokers),
		kafka.WithClientID("oms-producer"),
		kafka.WithKafkaVersion(version),  
	)

	kafkaLogger.Infof(" Kafka producer initialized with brokers: %v, version: %s", brokers, version)
}

func PublishOrderCreated(ctx context.Context, o *model.Order) {
	event := model.OrderCreated{
		OrderID:   o.ID,
		TenantID:  o.TenantID,
		SellerID:  o.SellerID,
		HubCode:     o.HubID,
		SKUCode:     o.SKUID,
		Quantity:  o.Quantity,
		CreatedAt: o.CreatedAt,
	}
	kafkaLogger.Infof(" Producer topic: %s", config.GetString(ctx, "kafka.producer_topic"))

	payload, err := pubsub.NewEventInBytes(event)

	if err != nil {
		kafkaLogger.Errorf(" Failed to marshal OrderCreated: %v", err)
		return
	}

	msg := &pubsub.Message{
		Topic: config.GetString(ctx, "kafka.producer_topic"),
		Key:   o.ID,
		Value: payload,
	}
	kafkaLogger.Infof(" About to publish to Kafka: topic=%s, key=%s, payload=%s", 
	msg.Topic, msg.Key, string(msg.Value))

	if err := producer.Publish(ctx, msg); err != nil {
		kafkaLogger.Errorf(" Kafka publish error: %v", err)
	} else {
		kafkaLogger.Infof(" Published order.created for OrderID: %s", o.ID)
	}
}
