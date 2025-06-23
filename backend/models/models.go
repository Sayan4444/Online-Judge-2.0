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
	Image 	  string    `json:"image"` 
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Contests []Contest `json:"contests" gorm:"many2many:contest_users;"`
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

	Submissions  []Submission `json:"submissions" gorm:"foreignKey:ProblemID"`
	Tests      []TestCase `json:"tests" gorm:"foreignKey:ProblemID"`
}

type Submission struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	ProblemID uuid.UUID `json:"problem_id" gorm:"not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	SubmittedAt time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	Result    string    `json:"result" gorm:"not null"` // e.g., "AC", "WA", "pending"
	Language  string    `json:"language" gorm:"not null"` // Programming language used for the submission
	SourceCode string    `json:"source_code" gorm:"not null"`
	Score	 int       `json:"score" gorm:"default:0"` 

	Problem   Problem  `json:"problem" gorm:"foreignKey:ProblemID"`
	User      User     `json:"user" gorm:"foreignKey:UserID"`
}

type TestCase struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	ProblemID uuid.UUID `json:"problem_id" gorm:"not null"`
	Input     string    `json:"input" gorm:"not null"` // Input for the test case
	Output    string    `json:"output" gorm:"not null"` // Expected output for the test case
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Problem   Problem  `json:"problem" gorm:"foreignKey:ProblemID"`
}