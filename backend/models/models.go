package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
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
}

type Submission struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	ProblemID uuid.UUID `json:"problem_id" gorm:"not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	SubmittedAt time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	Result    string    `json:"result" gorm:"not null"` // e.g., "accepted", "failed", "pending"
	Language  string    `json:"language" gorm:"not null"` // e.g., "C", "C++", "Python"
	SourceCode string    `json:"source_code" gorm:"not null"` // The actual code submitted by the user

	Problem   Problem  `json:"problem" gorm:"foreignKey:ProblemID"`
	User      User     `json:"user" gorm:"foreignKey:UserID"`
}