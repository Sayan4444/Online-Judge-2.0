package schema

import "github.com/google/uuid"

type RabbitMQPayload struct {
	SubmissionID   uuid.UUID `json:"submission_id"`
	ProblemID      uuid.UUID `json:"problem_id"`
	UserID         uuid.UUID `json:"user_id"`
	Language       string    `json:"language"`
	SourceCode     string    `json:"source_code"`
}

