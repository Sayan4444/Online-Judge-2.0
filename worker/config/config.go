/***********************************************************************
     Copyright (c) 2025 GNU/Linux Users' Group (NIT Durgapur)
************************************************************************/

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Load environment variables from env file
func LoadEnv() {
	// Load .env only if running locally
	if _, exists := os.LookupEnv("RUNNING_IN_DOCKER"); !exists {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found, relying on environment variables")
		}
	}
}

// retrieve env value from key
func GetEnv(key string) string {
	val := os.Getenv(key)

	return val
}

var RabbitMQConnection *amqp.Connection

func ConnectRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial(GetEnv("RABBITMQ_URL"))
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
		return nil, err
	}
	RabbitMQConnection = conn
	return conn, nil
}

func CreateRabbitMQChannel() (*amqp.Channel, error) {
	if RabbitMQConnection == nil {
		return nil, fmt.Errorf("RabbitMQ connection is not established")
	}
	
	ch, err := RabbitMQConnection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %s", err)
		return nil, err
	}
	return ch, nil
}

func CloseRabbitMQ() error {
	if RabbitMQConnection != nil {
		if err := RabbitMQConnection.Close(); err != nil {
			return err
		}
	}
	return nil
}
