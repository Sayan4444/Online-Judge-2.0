package rabbitmq

import (
	"OJ-backend/config"
	model "OJ-backend/models"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendSubmissionToQueue(rabbitmqPayload model.RabbitMQPayload, submissionID string) error {
	// create a channel
	ch, err := config.CreateRabbitMQChannel()
	if err != nil {
		log.Fatalf("Failed to create submit channel: %s", err)
	}
	defer ch.Close()
	// create a submission queue
	submissionQueue, err := ch.QueueDeclare("submissions", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare submission queue: %s", err)
	}
	// create a temporary queue for results
	const fiveHoursInMs = int32(5 * time.Hour / time.Millisecond)
	args := amqp.Table{
		"x-expires": fiveHoursInMs,
	}
	tempQueue, err := ch.QueueDeclare(submissionID, true, true, false, false, args)
	if err != nil {
		log.Fatalf("Failed to declare temporary results queue: %s", err)
	}
	body, err := json.Marshal(rabbitmqPayload)
	if err != nil {
		log.Fatalf("Failed to marshal rabbitmqPayload: %s", err)
	}
	// push data into that queue
	err = ch.Publish(
		"",                                 // exchange
		submissionQueue.Name, 				// routing key
		false,                              // mandatory
		false,                              // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			Timestamp:     time.Now(),
			DeliveryMode:  amqp.Persistent,
			ReplyTo:       tempQueue.Name,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish to temporary results queue: %s", err)
	}
	
	log.Printf("Submission sent to queue: %s", submissionQueue.Name)
	return nil
}
