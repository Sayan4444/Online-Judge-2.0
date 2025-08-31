package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"OJ-backend/config"
	handler "OJ-backend/controllers"
	model "OJ-backend/models"
)

// SSEClient represents a single SSE connection
type SSEClient struct {
	UserID         string
	SubmissionID   string
	ResponseWriter http.ResponseWriter
	Done           chan bool
	Created        time.Time
}

// SubmissionUpdate represents the data sent via SSE
type SubmissionUpdate struct {
	SubmissionID  string              `json:"submission_id"`
	Result        string              `json:"result"`
	Score         int                 `json:"score"`
	StdOutput     string              `json:"std_output"`
	StdError      string              `json:"std_error"`
	CompileOutput string              `json:"compile_output"`
	ExitSignal    int                 `json:"exit_signal"`
	ExitCode      int                 `json:"exit_code"`
	Time          string              `json:"time"`
	Memory        string              `json:"memory"`
	Message       string              `json:"message,omitempty"`
	WrongAnswers  []WrongAnswer `json:"wrong_answers,omitempty"`
	Status        string              `json:"status"` // "completed", "error", etc.
}

// sendSSEMessage sends a formatted SSE message to a client
func sendSSEMessage(client *SSEClient, update SubmissionUpdate) error {
	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %v", err)
	}

	message := fmt.Sprintf("data: %s\n\n", data)

	if _, err := client.ResponseWriter.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write SSE message: %v", err)
	}

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
		Done:           make(chan bool, 1),
		Created:        time.Now(),
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
		if err := handleSubmissionCallback(data, submissionID, userID, client); err != nil {
			log.Printf("Failed to handle submission callback: %v", err)
			errorUpdate := SubmissionUpdate{
				SubmissionID: submissionID,
				Status:       "error",
				Message:      fmt.Sprintf("Failed to process result: %v", err),
			}
			sendSSEMessage(client, errorUpdate)
		}
	}()

	// Wait for completion or client disconnect
	select {
	case <-client.Done:
		log.Printf("SSE connection closed for user %s, submission %s", userID, submissionID)
	case <-c.Request().Context().Done():
		log.Printf("SSE connection cancelled for user %s, submission %s", userID, submissionID)
	}

	return nil
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
	case <-time.After(30 * time.Second): // optional timeout
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
	Stdout     string `json:"stdout"`
}

func handleSubmissionCallback(data []byte, submissionID string, userID string, client *SSEClient) error {
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

	if err := db.First(&submission, "id = ?", submissionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("submission not found")
		}
		return fmt.Errorf("database error: %v", err)
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
