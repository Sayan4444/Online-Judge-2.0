package schema

import "github.com/google/uuid"

type RabbitMQPayload struct {
	SubmissionID   uuid.UUID `json:"submission_id"`
	ProblemID      uuid.UUID `json:"problem_id"`
	UserID         uuid.UUID `json:"user_id"`
	Language       string    `json:"language"`
	SourceCode     string    `json:"source_code"`
	SourceFileName string    `json:"source_file_name"`
	Status         string    `json:"status"`
	Score          int       `json:"score"`
	TimeLimit      int       `json:"time_limit"`
	WallTimeLimit  int       `json:"wall_time_limit"`
	MemoryLimit    int       `json:"memory_limit"`
	StackLimit     int       `json:"stack_limit"`
	OutputLimit    int       `json:"output_limit"`
	StdIn          string    `json:"stdin"`
	StdOut         string    `json:"stdout"`
	CompileCmd     string    `json:"compile_cmd"`
	RunCmd         string    `json:"run_cmd"`
	CallBackURL    string    `json:"callback_url"`
}
