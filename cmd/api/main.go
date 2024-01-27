package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"L0-wb-intern-task/internal/config"
	"L0-wb-intern-task/internal/http-server/handlers"
	"L0-wb-intern-task/internal/nats/consumer"
	"L0-wb-intern-task/internal/nats/publisher"
	"L0-wb-intern-task/internal/storage/cache"
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

	log.Println("Init cache")
	cache := cache.NewCache()
	if err := cache.AddAllOrdersToCache(db); err != nil {
		log.Fatal("Error adding orders to cache:", err)
	}

	log.Println("Starting nats")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}

	sub, err := consumer.SubscribeAndConsume(db, nc, config.Nats.StreamName)
	if err != nil {
		log.Fatal("Error subscribing and consuming orders:", err)
	}

	publishFlag := flag.Bool("p", false, "Enable publishing orders")

	flag.Parse()

	if *publishFlag {
		publisher.PublishOrders(nc)
		log.Println("Publishing orders...")
	} else {
		log.Println("Publish flag is not set. Skipping publishing orders.")
	}
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: mux,
	}

	mux.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetOrderHandler(w, r, cache)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.HomeHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Server is running on :%d...\n", 8080)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error:", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for range stop {
		log.Println("Received signal. Exiting...")
		sub.Unsubscribe()
		nc.Close()
		return
	}
}
