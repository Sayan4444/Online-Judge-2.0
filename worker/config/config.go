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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

var DB *gorm.DB
var RabbitMQConnection *amqp.Connection

// Connect to database
func ConnectDB() (*gorm.DB, error) {
	dsn := GetEnv("DSN_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	DB = db
	return db, err
}

func CloseDB() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err != nil {
            return err
        }
        return sqlDB.Close()
    }
    return nil
}

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
