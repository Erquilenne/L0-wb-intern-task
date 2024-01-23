package publisher

import (
	"L0-wb-intern-task/internal/config"
	"L0-wb-intern-task/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
)

func PublishReviews(js nats.JetStreamContext) {
	orders, err := getOrders()
	if err != nil {
		log.Println(err)
		return
	}

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	for _, oneOrder := range orders {

		// create random message intervals to slow down
		r := rand.Intn(1500)
		time.Sleep(time.Duration(r) * time.Millisecond)

		reviewString, err := json.Marshal(oneOrder)
		if err != nil {
			log.Println(err)
			continue
		}

		// publish to REVIEWS.rateGiven subject
		_, err = js.Publish(config.Nats.SubjectNameOrderCreated, reviewString)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Publisher  =>  Message: %d\n", oneOrder.OrderID)
		}
	}
}

func getOrders() ([]models.Order, error) {
	rawOrder, _ := ioutil.ReadFile("./reviews.json")
	var orderObj []models.Order
	err := json.Unmarshal(rawOrder, &orderObj)

	return orderObj, err
}
