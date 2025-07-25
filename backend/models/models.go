package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	OauthID   string    `json:"oauth_id" gorm:"unique"`
	Provider  string    `json:"provider"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Contests    []Contest    `json:"contests" gorm:"many2many:contest_users;"`
	Submissions []Submission `json:"submissions" gorm:"foreignKey:UserID"`
}

type Contest struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	Problems []Problem `json:"problems" gorm:"foreignKey:ContestID"`
	Users    []User    `json:"users" gorm:"many2many:contest_users;"`
}

type Problem struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	ContestID   uuid.UUID `json:"contest_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	Submissions []Submission `json:"submissions" gorm:"foreignKey:ProblemID"`
	Tests       []TestCase   `json:"tests" gorm:"foreignKey:ProblemID"`
}

type Submission struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey"`
	ProblemID      uuid.UUID `json:"problem_id" gorm:"not null"`
	UserID         uuid.UUID `json:"user_id" gorm:"not null"`
	ContestID      uuid.UUID `json:"contest_id" gorm:"not null"`
	SubmittedAt    time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	Result         string    `json:"result" gorm:"not null"`   // e.g., "AC", "WA", "pending"
	Language       string    `json:"language" gorm:"not null"` // Programming language used for the submission
	SourceCode     string    `json:"source_code" gorm:"not null"`
	Score          int       `json:"score" gorm:"default:0"`
	StdInput       string    `json:"std_input"`
	ExpectedOutput string    `json:"expected_output"`
	StdOutput      string    `json:"std_output"`
	StdError       string    `json:"std_error"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	CompileOutput  string    `json:"compile_output"` // Output from the compilation process
	ExitSignal     int       `json:"exit_signal"`    // Exit signal from the execution of the code
	ExitCode       int       `json:"exit_code"`      // Exit code from the execution of the code
	CallbackURL    string    `json:"callback_url"`   // URL to send the result of the submission

	Problem Problem `json:"problem" gorm:"foreignKey:ProblemID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
}

type TestCase struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	ProblemID uuid.UUID `json:"problem_id" gorm:"not null"`
	Input     string    `json:"input" gorm:"not null"`  // Input for the test case
	Output    string    `json:"output" gorm:"not null"` // Expected output for the test case
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Problem Problem `json:"problem" gorm:"foreignKey:ProblemID"`
}

type Language struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	Name           string `json:"name" gorm:"not null"`            // Name of the programming language
	CompileCommand string `json:"compile_command" gorm:"not null"` // Command to compile the code
	RunCommand     string `json:"run_command" gorm:"not null"`     // Command to run
	TimeLimit      int    `json:"time_limit" gorm:"not null"`      // Time limit for the submission in milliseconds
	MemoryLimit    int    `json:"memory_limit"`    // Memory limit for the submission in MB
	WallLimit      int    `json:"wall_limit"`      // Wall time limit for the submission in seconds
	StackLimit     int    `json:"stack_limit"`     // Stack limit for the submission in MB
	OutputLimit    int    `json:"output_limit"`    // Output limit for the submission in MB
	SrcFile        string `json:"src_file" gorm:"not null"`        // Source file name for the submission
}

type LeaderboardEntry struct {
	UserID          uuid.UUID `json:"user_id"`
	Username        string    `json:"username"`
	TotalScore      int       `json:"total_score"`
	FirstSubmission time.Time `json:"first_submission"`
}

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
