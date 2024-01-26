package consumer

import (
	"encoding/json"
	"log"
	"time"

	"L0-wb-intern-task/internal/storage/pgsql"
	"L0-wb-intern-task/models"

	"github.com/nats-io/nats.go"
)

// ConsumeOrder consumes an order from JSON and saves it to the database
func consumeOrder(jsonOrder []byte, db *pgsql.Database) error {
	var order models.Order

	// Unmarshal JSON into the Order struct
	err := json.Unmarshal(jsonOrder, &order)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	// Save the order to the database
	log.Println(
		"OrderUID:", order.OrderUID,
		"\nTrackNumber:", order.TrackNumber,
		"\nEntry:", order.Entry,
		"\nDeliveryName:", order.Delivery.Name,
		"\nDeliveryPhone:", order.Delivery.Phone,
		"\nDeliveryZip:", order.Delivery.Zip,
		"\nDeliveryCity:", order.Delivery.City,
		"\nDeliveryAddress:", order.Delivery.Address,
		"\nDeliveryRegion:", order.Delivery.Region,
		"\nDeliveryEmail:", order.Delivery.Email,
		"\nTransaction:", order.Payment.Transaction,
		"\nRequestID:", order.Payment.RequestID,
		"\nCurrency:", order.Payment.Currency,
		"\nProvider:", order.Payment.Provider,
		"\nAmount:", order.Payment.Amount,
		"\nPaymentDT:", time.Unix(order.Payment.PaymentDT, 0),
		"\nBank:", order.Payment.Bank,
		"\nDeliveryCost:", order.Payment.DeliveryCost,
		"\nGoodsTotal:", order.Payment.GoodsTotal,
		"\nCustomFee:", order.Payment.CustomFee,
		"\nLocale:", order.Locale,
		"\nInternalSignature:", order.InternalSignature,
		"\nCustomerID:", order.CustomerID,
		"\nDeliveryService:", order.DeliveryService,
		"\nShardKey:", order.ShardKey,
		"\nSMID:", order.SMID,
		"\nDateCreated:", order.DateCreated,
		"\nOofShard:", order.OofShard,
	)
	err = db.SaveOrder(order)
	if err != nil {
		log.Println("Error saving order to the database:", err)
		return err
	}

	log.Println("Order saved successfully.")

	return nil
}

// SubscribeAndConsume subscribes to the NATS stream and consumes incoming orders
func SubscribeAndConsume(db *pgsql.Database, nc *nats.Conn, streamName string) (*nats.Subscription, error) {
	// Subscribe to the NATS stream
	sub, err := nc.Subscribe(streamName, func(msg *nats.Msg) {
		// Consume the order from the message payload
		err := consumeOrder(msg.Data, db)
		if err != nil {
			log.Println("Error consuming order:", err)
		}
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS stream:", err)
		return nil, err
	}

	return sub, nil
}
