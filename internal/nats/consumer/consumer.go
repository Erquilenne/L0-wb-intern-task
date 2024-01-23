package consumer

import (
	"L0-wb-intern-task/internal/config"
	"L0-wb-intern-task/models"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

func ConsumeOrders(js nats.JetStreamContext) {
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	_, err := js.Subscribe(config.Nats.SubjectNameOrderCreated, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			log.Println("Unable to Ack", err)
			return
		}

		var order models.Order
		err = json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Consumer  =>  Subject: %s  -  ID: %s  -  CustomerID: %s", m.Subject, order.OrderID, order.CustomerID)
	})

	if err != nil {
		log.Println("Subscribe failed")
		return
	}
}
