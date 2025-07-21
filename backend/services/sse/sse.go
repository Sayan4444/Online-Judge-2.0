package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// SSEManager manages all SSE connections
type SSEManager struct {
	clients    map[string]map[string]*SSEClient // userID -> submissionID -> client
	clientsMux sync.RWMutex
}

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
	SubmissionID  string `json:"submission_id"`
	Result        string `json:"result"`
	Score         int    `json:"score"`
	StdOutput     string `json:"std_output"`
	StdError      string `json:"std_error"`
	CompileOutput string `json:"compile_output"`
	ExitSignal    int    `json:"exit_signal"`
	ExitCode      int    `json:"exit_code"`
	Time          string `json:"time"`
	Memory        string `json:"memory"`
	Message       string `json:"message"`
	Status        string `json:"status"` // "completed", "error", etc.
}

var GlobalSSEManager *SSEManager

func init() {
	GlobalSSEManager = &SSEManager{
		clients: make(map[string]map[string]*SSEClient),
	}

	// Start cleanup routine for expired connections
	go GlobalSSEManager.cleanupExpiredConnections()
}

// AddClient adds a new SSE client connection
func (m *SSEManager) AddClient(userID, submissionID string, w http.ResponseWriter) *SSEClient {
	m.clientsMux.Lock()
	defer m.clientsMux.Unlock()

	client := &SSEClient{
		UserID:         userID,
		SubmissionID:   submissionID,
		ResponseWriter: w,
		Done:           make(chan bool, 1),
		Created:        time.Now(),
	}

	if m.clients[userID] == nil {
		m.clients[userID] = make(map[string]*SSEClient)
	}

	m.clients[userID][submissionID] = client
	log.Printf("Added SSE client for user %s, submission %s", userID, submissionID)

	return client
}

// RemoveClient removes an SSE client connection
func (m *SSEManager) RemoveClient(userID, submissionID string) {
	m.clientsMux.Lock()
	defer m.clientsMux.Unlock()

	if userClients, exists := m.clients[userID]; exists {
		if client, exists := userClients[submissionID]; exists {
			close(client.Done)
			delete(userClients, submissionID)

			if len(userClients) == 0 {
				delete(m.clients, userID)
			}

			log.Printf("Removed SSE client for user %s, submission %s", userID, submissionID)
		}
	}
}

// BroadcastToUser sends an update to a specific user's submission
func (m *SSEManager) BroadcastToUser(userID, submissionID string, update SubmissionUpdate) {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	if userClients, exists := m.clients[userID]; exists {
		if client, exists := userClients[submissionID]; exists {
			if err := m.sendSSEMessage(client, update); err != nil {
				log.Printf("Failed to send SSE message to user %s, submission %s: %v", userID, submissionID, err)
				// Remove the client if sending fails
				go m.RemoveClient(userID, submissionID)
			} else {
				log.Printf("Sent SSE update to user %s, submission %s", userID, submissionID)
				// Close connection after sending the update
				go func() {
					time.Sleep(100 * time.Millisecond) // Small delay to ensure message is sent
					m.RemoveClient(userID, submissionID)
				}()
			}
		}
	}
}

// sendSSEMessage sends a formatted SSE message to a client
func (m *SSEManager) sendSSEMessage(client *SSEClient, update SubmissionUpdate) error {
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

// cleanupExpiredConnections removes connections that have been open too long
func (m *SSEManager) cleanupExpiredConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.clientsMux.Lock()
		for userID, userClients := range m.clients {
			for submissionID, client := range userClients {
				// Remove connections older than 5 minutes
				if time.Since(client.Created) > 5*time.Minute {
					log.Printf("Cleaning up expired SSE connection for user %s, submission %s", userID, submissionID)
					close(client.Done)
					delete(userClients, submissionID)
				}
			}
			if len(userClients) == 0 {
				delete(m.clients, userID)
			}
		}
		m.clientsMux.Unlock()
	}
}

// HandleSSEConnection handles incoming SSE connection requests
func HandleSSEConnection(c echo.Context) error {
	userID := c.Param("user_id")
	submissionID := c.Param("submission_id")

	if userID == "" || submissionID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id and submission_id are required"})
	}

	// Set SSE headers
	w := c.Response().Writer
	h := c.Response().Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Headers", "Cache-Control")

	// Add client to manager
	client := GlobalSSEManager.AddClient(userID, submissionID, w)

	// Send initial connection message
	initialUpdate := SubmissionUpdate{
		SubmissionID: submissionID,
		Status:       "connected",
		Message:      "Connected to submission updates",
	}

	if err := GlobalSSEManager.sendSSEMessage(client, initialUpdate); err != nil {
		log.Printf("Failed to send initial SSE message: %v", err)
		GlobalSSEManager.RemoveClient(userID, submissionID)
		return err
	}

	// Wait for completion or client disconnect
	select {
	case <-client.Done:
		log.Printf("SSE connection closed for user %s, submission %s", userID, submissionID)
	case <-c.Request().Context().Done():
		log.Printf("SSE connection cancelled for user %s, submission %s", userID, submissionID)
		GlobalSSEManager.RemoveClient(userID, submissionID)
	}

	return nil
}
