package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"OJ-backend/config"
	handler "OJ-backend/controllers"
	model "OJ-backend/models"
)

// SSEClient represents a single SSE connection
type SSEClient struct {
	UserID         string
	SubmissionID   string
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Done           chan bool
	Created        time.Time
	RemoteAddr     string
	UserAgent      string
	LastActivity   time.Time
	MessagesSent   int
	BytesSent      int64
}

// SubmissionUpdate represents the data sent via SSE
type SubmissionUpdate struct {
	SubmissionID  string        `json:"submission_id"`
	Result        string        `json:"result"`
	Score         int           `json:"score"`
	StdOutput     string        `json:"std_output"`
	StdError      string        `json:"std_error"`
	CompileOutput string        `json:"compile_output"`
	ExitSignal    int           `json:"exit_signal"`
	ExitCode      int           `json:"exit_code"`
	Time          string        `json:"time"`
	Memory        string        `json:"memory"`
	Message       string        `json:"message,omitempty"`
	WrongAnswers  []WrongAnswer `json:"wrong_answers,omitempty"`
	Status        string        `json:"status"` // "completed", "error", etc.
}

// sendSSEMessage sends a formatted SSE message to a client
func sendSSEMessage(client *SSEClient, update SubmissionUpdate) error {
	if client.Request.Context().Err() != nil {
		// Context is done (e.g., canceled), so the client has disconnected.
		// Return the context's error to signal that we should stop sending.
		return client.Request.Context().Err()
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %v", err)
	}

	message := fmt.Sprintf("data: %s\n\n", data)

	bytesWritten, err := client.ResponseWriter.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write SSE message: %v", err)
	}

	// Update client tracking
	client.MessagesSent++
	client.BytesSent += int64(bytesWritten)
	client.LastActivity = time.Now()

	if flusher, ok := client.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// HandleSSEConnection handles incoming SSE connection requests
func HandleSSEConnection(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*handler.Claims)
	userID := claims.UserID
	submissionID := c.Param("submission_id")

	if userID == "" || submissionID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id and submission_id are required"})
	}
	request := c.Request()
	w := c.Response().Writer
	h := c.Response().Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Headers", "Cache-Control")

	client := &SSEClient{
		UserID:         userID,
		SubmissionID:   submissionID,
		ResponseWriter: w,
		Request:        request,
		Done:           make(chan bool, 1),
		Created:        time.Now(),
		RemoteAddr:     request.RemoteAddr,
		UserAgent:      request.Header.Get("User-Agent"),
		LastActivity:   time.Now(),
		MessagesSent:   0,
		BytesSent:      0,
	}

	// Send initial connection message
	initialUpdate := SubmissionUpdate{
		SubmissionID: submissionID,
		Status:       "connected",
		Message:      "Connected to submission updates",
	}

	if err := sendSSEMessage(client, initialUpdate); err != nil {
		log.Printf("Failed to send initial SSE message: %v", err)
		return err
	}

	// Start goroutine to wait for queue data and handle it
	go func() {
		defer func() {
			client.Done <- true
		}()

		// consumes result from queue and returns the data
		log.Printf("Consuming consume result from queue: %v", submissionID)
		data, err := consumeResult(submissionID)
		if err != nil {
			log.Printf("Failed to consume result from queue: %v", err)
			errorUpdate := SubmissionUpdate{
				SubmissionID: submissionID,
				Status:       "error",
				Message:      fmt.Sprintf("Failed to get result: %v", err),
			}
			sendSSEMessage(client, errorUpdate)
			return
		}

		// Handle the submission result
		if err := handleSubmissionCallback(c,data, submissionID, userID, client); err != nil {
			log.Printf("Failed to handle submission callback: %v", err)
			errorUpdate := SubmissionUpdate{
				SubmissionID: submissionID,
				Status:       "error",
				Message:      fmt.Sprintf("Failed to process result: %v", err),
			}
			sendSSEMessage(client, errorUpdate)
		}
	}()

	keepAliveTicker := time.NewTicker(30 * time.Second)
	defer keepAliveTicker.Stop()

	// Wait for completion or client disconnect
	for {
		select {
		case <-client.Done:
			duration := time.Since(client.Created)
			log.Printf("SSE connection completed normally - User: %s, Submission: %s, Duration: %v, Messages: %d, Bytes: %d, RemoteAddr: %s", 
				userID, submissionID, duration, client.MessagesSent, client.BytesSent, client.RemoteAddr)
			return nil
		case <-c.Request().Context().Done():
			contextErr := c.Request().Context().Err()
			duration := time.Since(client.Created)
			log.Printf("SSE connection cancelled - User: %s, Submission: %s, Duration: %v, Reason: %v, Messages: %d, Bytes: %d, RemoteAddr: %s, UserAgent: %s", 
				userID, submissionID, duration, contextErr, client.MessagesSent, client.BytesSent, client.RemoteAddr, client.UserAgent)
			return nil
		case <-keepAliveTicker.C:
			keepAliveUpdate := SubmissionUpdate{
				SubmissionID: submissionID,
				Status:       "keep-alive",
				Message:      "Keep-alive ping",
			}
			if err := sendSSEMessage(client, keepAliveUpdate); err != nil {
				log.Printf("Failed to send keep-alive message for user %s, submission %s: %v", userID, submissionID, err)
				return nil
			}
			log.Printf("Sent keep-alive message - User: %s, Submission: %s", userID, submissionID)
		}
	}
}

func consumeResult(submissionID string) ([]byte, error) {
	ch, err := config.CreateRabbitMQChannel()
	if err != nil {
		log.Fatalf("Failed to create submit channel: %s", err)
	}
	defer ch.Close()
	msgs, err := ch.Consume(submissionID, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	select {
	case d := <-msgs:
		d.Ack(false)
		return d.Body, nil
	case <-time.After(1 * time.Hour): // optional timeout
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

type receivedPayload struct {
	SubmissionID  uuid.UUID     `json:"submission_id"`
	Score         int           `json:"score"`
	JudgeResponse JudgeResponse `json:"judge_response"`
}

type JudgeResponse struct {
	Stderr        string        `json:"stderr"`
	Time          string        `json:"time"`
	Memory        string        `json:"memory"`
	ExitCode      string        `json:"exit_code"`
	Result        string        `json:"result"`
	CompileOutput string        `json:"compile_output"`
	WrongAnswers  []WrongAnswer `json:"wrong_answers"`
}

type WrongAnswer struct {
	TestCaseID uuid.UUID `json:"test_case_id"`
	Stdout     string    `json:"stdout"`
}

func handleSubmissionCallback(c echo.Context, data []byte, submissionID string, userID string, client *SSEClient) error {
	// Parse the callback payload from queue data
	var payload receivedPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Failed to unmarshal judge response: %v", err)
		return fmt.Errorf("failed to parse judge response: %v", err)
	}

	judgeResponse := payload.JudgeResponse

	// Log all fields of judgeResponse for debugging
	log.Printf("Judge Response Details:")
	log.Printf("  Stderr: %s", judgeResponse.Stderr)
	log.Printf("  Time: %s", judgeResponse.Time)
	log.Printf("  Memory: %s", judgeResponse.Memory)
	log.Printf("  ExitCode: %s", judgeResponse.ExitCode)
	log.Printf("  Result: %s", judgeResponse.Result)
	log.Printf("  CompileOutput: %s", judgeResponse.CompileOutput)

	db := config.DB
	var submission model.Submission

	redisClient := config.RedisClient
	ctx := c.Request().Context()

	// Get submission data from Redis
	submissionJSON, err := redisClient.Get(ctx, "submission:"+submissionID).Result()
	if err != nil {
		log.Printf("Failed to get submission %s from Redis: %v", submissionID, err)
		return fmt.Errorf("failed to retrieve submission from cache: %v", err)
	}

	if err := json.Unmarshal([]byte(submissionJSON), &submission); err != nil {
		log.Printf("Failed to unmarshal submission %s from Redis: %v", submissionID, err)
		return fmt.Errorf("failed to parse submission data: %v", err)
	}

	// Update submission with results
	submission.Result = judgeResponse.Result
	submission.StdError = judgeResponse.Stderr
	submission.CompileOutput = judgeResponse.CompileOutput
	if len(judgeResponse.WrongAnswers) > 0 {
		submission.WrongTestCase = judgeResponse.WrongAnswers[0].TestCaseID
		submission.StdOutput = judgeResponse.WrongAnswers[0].Stdout
	}

	if exitCode, err := strconv.Atoi(judgeResponse.ExitCode); err == nil {
		submission.ExitCode = exitCode
	} else {
		submission.ExitCode = 0
		log.Printf("Failed to convert exit code '%s' to int: %v", judgeResponse.ExitCode, err)
	}

	if err := db.Save(&submission).Error; err != nil {
		return fmt.Errorf("failed to update submission: %v", err)
	}

	// Create SSE update and send to client
	sseUpdate := SubmissionUpdate{
		SubmissionID:  submissionID,
		Result:        judgeResponse.Result,
		Score:         submission.Score, // Use existing score from DB
		StdError:      judgeResponse.Stderr,
		CompileOutput: judgeResponse.CompileOutput,
		ExitCode:      submission.ExitCode,
		Time:          judgeResponse.Time,
		Memory:        judgeResponse.Memory,
		WrongAnswers:  judgeResponse.WrongAnswers,
		Status:        "completed",
	}

	// Send SSE message to client
	if err := sendSSEMessage(client, sseUpdate); err != nil {
		log.Printf("Failed to send SSE message: %v", err)
		return fmt.Errorf("failed to send SSE message: %v", err)
	}

	log.Printf("Successfully processed submission %s with result %s", submissionID, judgeResponse.Result)
	return nil
}
