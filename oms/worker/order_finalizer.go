package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dhruv/oms/client"
	"github.com/dhruv/oms/model"
)

type OrderCreatedHandler struct{}

func (h *OrderCreatedHandler) Process(ctx context.Context, msg *pubsub.Message) error {
	var event model.OrderCreated
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.DefaultLogger().Errorf(" Failed to unmarshal OrderCreated: %v", err)
		return err
	}

	logger := log.DefaultLogger()
	logger.Infof(" Processing order.created for OrderID: %s", event.OrderID)

	baseURL := config.GetString(ctx, "ims.base_url")
	timeout := config.GetDuration(ctx, "ims.timeout")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Call IMS to check inventory
	inventory, err := client.FetchInventory(ctxWithTimeout, baseURL, event.TenantID, event.SellerID, event.HubCode, event.SKUCode)
	if err != nil {
		logger.Errorf(" IMS fetch inventory failed: %v", err)
		return err
	}

	if inventory.Quantity >= event.Quantity {
		// Reduce inventory
		if err := client.ConsumeInventory(ctxWithTimeout, baseURL, event.TenantID, event.SellerID, event.HubCode, event.SKUCode, event.Quantity); err != nil {
			logger.Errorf(" IMS consume inventory failed: %v", err)
			return err
		}

		// Update order status to new_order
		if err := client.UpdateOrderStatus(ctxWithTimeout, baseURL, client.UpdateOrderStatusRequest{
			OrderID: event.OrderID,
			Status:  "new_order",
		}); err != nil {
			logger.Errorf(" Failed to update order status: %v", err)
			return err
		}
		logger.Infof(" Order %s finalized as new_order", event.OrderID)

		// 🔔 NEW: Trigger webhook

		// client.NotifyWebhooks(ctx, event.TenantID, "order.updated", map[string]interface{}{
		// 	"order_id": event.OrderID,
		// 	"status":   "new_order",
		// })
		///////////////////
		coll, err := client.GetOrdersCollection(ctx)
		if err != nil {
			logger.Errorf(" Failed to get orders collection: %v", err)
			return err
		}

		var fullOrder model.Order
		if err := coll.FindOne(ctx, bson.M{"_id": event.OrderID}).Decode(&fullOrder); err != nil {
			logger.Errorf(" Failed to load order for webhook: %v", err)
			return err
		}

		// Trigger webhook with full order
		client.NotifyWebhooks(ctx, fullOrder.TenantID, "order.updated", fullOrder)

		//////////////////////////
	} else {
		// Not enough stock, keep on_hold
		if err := client.UpdateOrderStatus(ctxWithTimeout, baseURL, client.UpdateOrderStatusRequest{
			OrderID: event.OrderID,
			Status:  "on_hold",
		}); err != nil {
			logger.Errorf(" Failed to update order status: %v", err)
			return err
		}
		logger.Warnf(" Order %s kept on_hold due to insufficient inventory", event.OrderID)
	}

	return nil
}

func StartOrderFinalizer(ctx context.Context) {
	brokers := config.GetStringSlice(ctx, "kafka.brokers")
	groupID := config.GetString(ctx, "kafka.consumer_group")
	version := config.GetString(ctx, "kafka.version")
	clientID := config.GetString(ctx, "kafka.producer_topic")

	log.DefaultLogger().Infof("Kafka config: brokers=%v version=%s", brokers, version)

	consumer := kafka.NewConsumer(
		kafka.WithBrokers(brokers),
		kafka.WithConsumerGroup(groupID),
		kafka.WithClientID(clientID),
		kafka.WithKafkaVersion(version),
		kafka.WithRetryInterval(time.Second),
	)

	log.DefaultLogger().Infof(" Consumer subscribing to topic: order.created")

	handler := &OrderCreatedHandler{
		
	}
	consumer.RegisterHandler("order.created", handler)

	go consumer.Subscribe(ctx)
}
