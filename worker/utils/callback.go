package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// CallbackPayload represents the data sent back to the server
type CallbackPayload struct {
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
}

// generateHMAC generates HMAC-SHA256 signature for the payload
func generateHMAC(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// SendCallback sends HMAC-authenticated callback to the server
func SendCallback(callbackURL string, payload CallbackPayload, webhookSecret string) error {
	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Generate HMAC signature
	signature := generateHMAC(jsonPayload, webhookSecret)

	// Create HTTP request
	req, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OJ-Signature", fmt.Sprintf("sha256=%s", signature))
	req.Header.Set("User-Agent", "OJ-Worker/1.0")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send callback: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("callback failed with status: %d", resp.StatusCode)
	}

	return nil
}
