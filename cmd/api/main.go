package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"L0-wb-intern-task/internal/config"
	"L0-wb-intern-task/internal/nats/consumer"
	"L0-wb-intern-task/internal/nats/publisher"
	"L0-wb-intern-task/internal/storage/pgsql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)

	db, err := pgsql.NewDatabase(connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	defer db.Close()

	db.MakeMigrations()

	log.Println("Starting nats")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}

	sub, err := consumer.SubscribeAndConsume(db, nc, config.Nats.StreamName)
	if err != nil {
		log.Fatal("Error subscribing and consuming orders:", err)
	}

	publisher.PublishOrders(nc)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			log.Println("Received signal. Exiting...")
			sub.Unsubscribe()
			nc.Close()
			return
		}
	}
}
