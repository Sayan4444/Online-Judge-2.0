package main

import (
	"OJ-Worker/config"
	isolatejob "OJ-Worker/isolateJob"
	"OJ-Worker/schema"
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"strconv"
	"sync" // Import sync for WaitGroup
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// processMessage processes a submission message
func processMessage(ctx context.Context, d amqp.Delivery, workerTag string, ch *amqp.Channel) {
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

	// Marshal the publish payload to JSON
	body, err := json.Marshal(publishPayload)
	if err != nil {
		log.Printf("%s: Failed to marshal publish payload: %v", workerTag, err)
		return
	}

	err = ch.Publish(
		"",              // exchange
		d.ReplyTo,       // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Printf("%s: Failed to publish result: %v", workerTag, err)
		
		// Check if channel is closed
		if ch.IsClosed() {
			log.Printf("%s: [CRITICAL] Channel is closed during publish operation!", workerTag)
		}
		return
	}

	// Log all fields of publishPayload.JudgeResponse
	log.Printf("%s: JudgeResponse fields for submission %s:", workerTag, submission.SubmissionID)
	log.Printf("  Stderr: %s", publishPayload.JudgeResponse.Stderr)
	log.Printf("  Time: %s", publishPayload.JudgeResponse.Time)
	log.Printf("  Memory: %s", publishPayload.JudgeResponse.Memory)
	log.Printf("  ExitCode: %s", publishPayload.JudgeResponse.ExitCode)
	log.Printf("  Result: %s", publishPayload.JudgeResponse.Result)
	log.Printf("  CompileOutput: %s", publishPayload.JudgeResponse.CompileOutput)
	log.Printf("  WrongAnswers: %+v", publishPayload.JudgeResponse.WrongAnswers)
}

func main() {
	// Load environment variables
	config.LoadEnv()

	conn, err := config.ConnectRabbitMQ()
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	
	// Add connection close monitoring
	connCloseChan := make(chan *amqp.Error)
	conn.NotifyClose(connCloseChan)
	
	go func() {
		closeErr := <-connCloseChan
		if closeErr != nil {
			log.Printf("[CRITICAL] RabbitMQ connection closed unexpectedly: %v", closeErr)
			log.Printf("  Error Code: %d", closeErr.Code)
			log.Printf("  Error Reason: %s", closeErr.Reason)
			log.Printf("  Server Initiated: %t", closeErr.Server)
			log.Printf("  Recoverable: %t", closeErr.Recover)
		} else {
			log.Printf("[INFO] RabbitMQ connection closed gracefully")
		}
	}()

	// Configure the number of concurrent workers from environment variables.
	numWorkers, err := strconv.Atoi(config.GetEnv("NUM_WORKERS"))
	if err != nil {
		numWorkers = 5
	}

	log.Printf("Worker Configuration:")
	log.Printf("  NUM_WORKERS: %d", numWorkers)
	log.Printf("  RabbitMQ Connection: %s", config.GetEnv("RABBITMQ_URL"))
	log.Printf("Starting %d workers", numWorkers)
	
	// create channel
	ch, err := config.CreateRabbitMQChannel()
	if err != nil {
		log.Fatalf("Failed to create submit channel: %s", err)
	}
	defer ch.Close()
	
	// Add channel close monitoring
	channelCloseChan := make(chan *amqp.Error)
	ch.NotifyClose(channelCloseChan)
	
	go func() {
		closeErr := <-channelCloseChan
		if closeErr != nil {
			log.Printf("[CRITICAL] RabbitMQ channel closed unexpectedly: %v", closeErr)
			log.Printf("  Error Code: %d", closeErr.Code)
			log.Printf("  Error Reason: %s", closeErr.Reason)
			log.Printf("  Server Initiated: %t", closeErr.Server)
			log.Printf("  Recoverable: %t", closeErr.Recover)
		} else {
			log.Printf("[INFO] RabbitMQ channel closed gracefully")
		}
	}()
	// create submission queue
	submissionQueue, err := ch.QueueDeclare("submissions", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare submission queue: %s", err)
	}
	err = ch.Qos(
		numWorkers,
		0,
		false,
	)
	failOnError(err, "Failed to set QoS")
	// get code from there
	// do the processing
	// send results to the temp queue
	// close the channel
	msgs, err := ch.Consume(
		submissionQueue.Name, // queue name
		"",                       // consumer: empty string for auto-generated consumer tag
		false,                    // auto-ack: false for manual acknowledgement
		false,                    // exclusive: false allows multiple consumers
		false,                    // no-local: false means consume messages published by this connection
		false,                    // no-wait: don't wait for server confirmation
		nil,                      // args
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

			log.Printf("%s: Started and ready to process messages", workerTag)

			for {
				select {
				case d, ok := <-msgs:
					if !ok {
						log.Printf("%s: Message channel closed. Checking if intentional shutdown...", workerTag)
						
						// Check if this is due to context cancellation (intentional shutdown)
						select {
						case <-ctx.Done():
							log.Printf("%s: Channel closed due to graceful shutdown", workerTag)
						default:
							log.Printf("%s: [UNEXPECTED] Channel closed WITHOUT shutdown signal!", workerTag)
						}
						return
					}
					processMessage(ctx, d, workerTag,ch)
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
