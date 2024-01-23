package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"L0-wb-intern-task/internal/config"
	"L0-wb-intern-task/internal/nats"
	"L0-wb-intern-task/internal/nats/consumer"
	"L0-wb-intern-task/internal/nats/publisher"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Error creating database driver instance:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Error creating migration instance:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations:", err)
	}

	fmt.Println("Migrations applied successfully!")
	log.Println("Starting nats")

	js, err := nats.JetStreamInit()
	if err != nil {
		log.Println(err)
		return
	}
	// Let's assume that publisher and consumer are services running on different servers.
	// So run publisher and consumer asynchronously to see how it works
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisher.PublishOrders(js)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.ConsumeOrders(js)
	}()

	wg.Wait()

	log.Println("Exit...")
}
