package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
	"encoding/json"
	model "OJ-backend/models"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
}


var RabbitMQClient *RabbitMQ

func NewRabbitMQ(queueName string) *RabbitMQ {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	// Declare a queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		QueueName:  queueName,
	}
}

func SendSubmissionToQueue(rabbitmqPayload model.RabbitMQPayload) error {
	if RabbitMQClient == nil {
		RabbitMQClient = NewRabbitMQ("submissions")
	}
	body, err := json.Marshal(rabbitmqPayload)
	if err != nil {
		return err
	}

	err = RabbitMQClient.Channel.Publish(
		"",              // exchange
		RabbitMQClient.QueueName, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Submission sent to queue: %s", RabbitMQClient.QueueName)
	return nil
}

func CloseRabbitMQ() {
	if RabbitMQClient != nil {
		if RabbitMQClient.Channel != nil {
			RabbitMQClient.Channel.Close()
		}
		if RabbitMQClient.Connection != nil {
			RabbitMQClient.Connection.Close()
		}
		log.Println("RabbitMQ connection closed")
	}
}