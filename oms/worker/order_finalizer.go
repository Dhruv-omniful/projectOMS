package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"

	// "github.com/dhruv/oms/client"
	"github.com/dhruv/oms/model"
)

type OrderCreatedHandler struct{}

func (h *OrderCreatedHandler) Process(ctx context.Context, msg *pubsub.Message) error {
	var event model.OrderCreated
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.DefaultLogger().Errorf("❌ Failed to unmarshal OrderCreated: %v", err)
		return err
	}

	log.DefaultLogger().Infof("📥 Processing order.created for OrderID: %s", event.OrderID)

	// Example: Check inventory, update Mongo, etc.
	// if err := client.CheckAndUpdateOrder(ctx, &event); err != nil {
	// 	log.DefaultLogger().Errorf("❌ Failed to finalize order: %v", err)
	// 	return err
	// }

	log.DefaultLogger().Infof("✅ Finalized order: %s", event.OrderID)
	return nil
}

func StartOrderFinalizer(ctx context.Context) {
	brokers := config.GetStringSlice(ctx, "kafka.brokers")
	groupID := config.GetString(ctx, "kafka.consumer_group")
	version := config.GetString(ctx, "kafka.version") // Add in config.yaml
	clientID := config.GetString(ctx, "kafka.producer_topic")
	log.DefaultLogger().Infof("Kafka config: brokers=%v version=%s", brokers, version)
	consumer := kafka.NewConsumer(
		kafka.WithBrokers(brokers),
		kafka.WithConsumerGroup(groupID),
		kafka.WithClientID(clientID),
		kafka.WithKafkaVersion(version),
		kafka.WithRetryInterval(time.Second),
	)
	log.DefaultLogger().Infof("✅ Consumer subscribing to topic: order.created")

	// Optional: Add interceptor if needed (e.g. NewRelic)
	// consumer.SetInterceptor(interceptor.NewRelicInterceptor())

	// Register handler
	handler := &OrderCreatedHandler{}
	consumer.RegisterHandler("order.created", handler)

	go consumer.Subscribe(ctx)
}
