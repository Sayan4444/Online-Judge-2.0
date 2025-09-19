package schema

import (
	model "OJ-Worker/models"

	"github.com/google/uuid"
)

type RabbitMQPayload struct {
	SubmissionID uuid.UUID        `json:"submission_id"`
	ProblemID    uuid.UUID        `json:"problem_id"`
	UserID       uuid.UUID        `json:"user_id"`
	SourceCode   string           `json:"source_code"`
	Language     model.Language  `json:"language"`
	TestCases    []model.TestCase `json:"test_cases"`
}
