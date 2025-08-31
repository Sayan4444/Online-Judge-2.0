package rabbitmq

import (
	"OJ-backend/config"
	model "OJ-backend/models"
	"encoding/json"
	"fmt"
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
	const fiveMinutesInMs = int32(5 * time.Minute / time.Millisecond)
	args := amqp.Table{
		"x-expires": fiveMinutesInMs,
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
	
	// only for checking results
	// consumeResult(ch, tempQueue)

	log.Printf("Submission sent to queue: %s", submissionQueue.Name)
	return nil
}

// only for checking results in dev
func consumeResult(ch *amqp.Channel, tempQueue amqp.Queue) ([]byte, error) {
	defer ch.Close()
	msgs, err := ch.Consume(tempQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	select {
	case d := <-msgs:
		log.Printf(" [x] Got reply: %s", d.Body)
		d.Ack(false)
		return d.Body, nil
	case <-time.After(30 * time.Second): // optional timeout
		return nil, fmt.Errorf("timeout waiting for response")
	}
}
