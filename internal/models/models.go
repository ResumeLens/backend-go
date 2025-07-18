package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type User struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string `gorm:"unique;not null"`
	PasswordHash   string `gorm:"not null"`
	Role           string `gorm:"not null"`
	OrganizationID string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Organization struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string `gorm:"unique;not null"`
	CreatedBy string
	CreatedAt time.Time
}

type Invite struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string `gorm:"not null"`
	OrganizationID string `gorm:"not null"`
	Role           string `gorm:"not null"`
	Token          string `gorm:"unique;not null"`
	Expiry         time.Time
	IsAccepted     bool
	CreatedAt      time.Time
}

type Candidate struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FullName  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Password  string `gorm:""`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CandidateApplication struct {
	CandidateID    string `gorm:"primaryKey;type:uuid"`
	OrganizationID string `gorm:"primaryKey;type:uuid"`
	JobID          string `gorm:"primaryKey;type:uuid"`

	ResumeGCSPath string          `gorm:"not null"`
	ParsedResume  string          `gorm:"type:jsonb"`
	Embedding     pgvector.Vector `gorm:"type:vector(384)"` // all-MiniLM-L6-v2 -- 384 dimensions
	Status        string          `gorm:"not null;default:'pending'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
