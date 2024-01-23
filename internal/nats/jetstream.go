package main

import (
	"L0-wb-intern-task/internal/config"
	"log"

	"github.com/nats-io/nats.go"
)

func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	// Create a stream if it does not exist
	err = CreateStream(js)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext) error {

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	stream, err := jetStream.StreamInfo(config.Nats.StreamName)

	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", config.Nats.StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     config.Nats.StreamName,
			Subjects: []string{config.Nats.StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
