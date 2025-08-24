package main

import (
	"OJ-Worker/config"
	"OJ-Worker/isolateJob"
	"OJ-Worker/schema"
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"strconv"
	"sync" // Import sync for WaitGroup
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// processMessage processes a submission message
func processMessage(ctx context.Context, d amqp091.Delivery, workerTag string) {
	log.Printf("%s: Processing submission: %s", workerTag, d.Body)

	// Parse submission from message
	var submission schema.RabbitMQPayload
	if err := json.Unmarshal(d.Body, &submission); err != nil {
		log.Printf("%s: Failed to parse submission: %v", workerTag, err)
		d.Nack(false, false) // Don't requeue malformed messages
		return
	}

	// Initialize response
	response := &schema.JudgeResponse{}

	// Process submission using isolate
	if err := isolatejob.ProcessSubmission(&submission, response, ctx); err != nil {
		log.Printf("%s: Failed to process submission %s: %v", workerTag, submission.SubmissionID, err)
		response.Result = schema.ResultSystemError
		response.Message = "Internal processing error"
	}

	// Calculate score based on result
	score := 0
	if response.Result == schema.ResultAccepted {
		score = 100 // Full score for accepted solution
	}

	// Prepare callback payload
	publishPayload := schema.PublishPayload{
		SubmissionID:  submission.SubmissionID,
		Score:         score,
		JudgeResponse: *response,
	}

	// Will use pub sub	

	// Log all fields of publishPayload.JudgeResponse
	log.Printf("%s: JudgeResponse fields for submission %s:", workerTag, submission.SubmissionID)
	log.Printf("  Stderr: %s", publishPayload.JudgeResponse.Stderr)
	log.Printf("  Time: %s", publishPayload.JudgeResponse.Time)
	log.Printf("  Memory: %s", publishPayload.JudgeResponse.Memory)
	log.Printf("  ExitSignal: %s", publishPayload.JudgeResponse.ExitSignal)
	log.Printf("  ExitCode: %s", publishPayload.JudgeResponse.ExitCode)
	log.Printf("  Message: %s", publishPayload.JudgeResponse.Message)
	log.Printf("  Result: %s", publishPayload.JudgeResponse.Result)
	log.Printf("  CompileOutput: %s", publishPayload.JudgeResponse.CompileOutput)
	log.Printf("  WrongAnswers: %+v", publishPayload.JudgeResponse.WrongAnswers)
}

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
	} else {
		log.Printf("Successfully connected to the database: %s", db.Name())
	}

	// Configure the number of concurrent workers from environment variables.
	numWorkers, err := strconv.Atoi(config.GetEnv("NUM_WORKERS"))
	if err != nil {
		numWorkers = 5
	}

	log.Printf("Starting %d workers", numWorkers)

	// Configure RabbitMQ connection from environment variables.
	amqpURI := config.GetEnv("RABBITMQ_URL")
	if amqpURI == "" {
		amqpURI = "amqp://guest:guest@localhost:5672" // Default fallback
	}

	// 1. Connect to RabbitMQ
	conn, err := amqp091.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func() {
		log.Println("[*] Closing RabbitMQ connection...")
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}()

	// 2. Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func() {
		log.Println("[*] Closing RabbitMQ channel...")
		if err := ch.Close(); err != nil {
			log.Printf("Error closing channel: %s", err)
		}
	}()

	// 3. Declare a durable queue
	q, err := ch.QueueDeclare(
		"submissions", // name
		true,          // durable - messages will survive broker restarts
		false,         // delete when unused
		false,         // exclusive - only accessible by this connection
		false,         // no-wait - don't wait for server confirmation
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 4. Set Quality of Service (QoS) for prefetch
	// This ensures that RabbitMQ will send at most `numWorkers` messages to this consumer
	// that have not yet been acknowledged, distributing them among your workers.
	err = ch.Qos(
		numWorkers, // prefetch count: send `numWorkers` unacknowledged messages
		0,          // prefetch size: 0 means no limit on message size
		false,      // global: false means QoS applies per consumer
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer: empty string for auto-generated consumer tag
		false,  // auto-ack: false for manual acknowledgement
		false,  // exclusive: false allows multiple consumers
		false,  // no-local: false means consume messages published by this connection
		false,  // no-wait: don't wait for server confirmation
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM) //Graceful shutdown signal
	defer stop()

	var wg sync.WaitGroup

	log.Printf(" [*] Starting %d workers. Waiting for messages. To exit, press CTRL+C", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			workerTag := "worker-" + strconv.Itoa(workerID)

			for {
				select {
				case d, ok := <-msgs:
					if !ok {
						log.Printf("%s: Message channel closed. Worker exiting.", workerTag)
						return
					}
					processMessage(ctx, d, workerTag)
					d.Ack(false) // Acknowledge the message completed
				case <-ctx.Done():
					log.Printf("%s: Application context cancelled. Worker exiting gracefully.", workerTag)
					return
				}
			}
		}(i)
	}

	<-ctx.Done()
	log.Println("\n[*] Received shutdown signal. Initiating graceful shutdown...")

	waitTimeout := 15 * time.Second //Timeout for the workers to finish when graceful shutdown is initiated else forcefully shut down
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[*] All message processing goroutines have completed.")
	case <-time.After(waitTimeout):
		log.Printf("[!] Timeout (%s) waiting for all message processing goroutines. Forcefully shutting down.", waitTimeout)
	}

	log.Println("[*] Application exiting.")
}
