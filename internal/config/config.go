package config

import (
	"encoding/json"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type NatsConfig struct {
	StreamName               string `json:"stream_name"`
	StreamSubjects           string `json:"stream_subjects"`
	SubjectNameOrderCreated  string `json:"subject_name_order_created"`
	SubjectNameOrderAnswered string `json:"subject_name_order_answered"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Nats     NatsConfig     `json:"nats"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
