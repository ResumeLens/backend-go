package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type User struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string `gorm:"unique;not null"`
	PasswordHash   string `gorm:"not null"`
	RoleID         string `gorm:"not null"`
	OrganizationID string `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type Organization struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `gorm:"unique;not null"`
	CreatedByID *string `gorm:"type:uuid;column:created_by"`
	CreatedAt   time.Time
}

type Invite struct {
	ID             string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID         *string `gorm:"type:uuid"`
	Email          string  `gorm:"not null"`
	OrganizationID string  `gorm:"not null"`
	RoleID         string  `gorm:"not null"`
	Token          string  `gorm:"unique;not null"`
	Expiry         time.Time
	IsAccepted     bool
	CreatedAt      time.Time
}

type Candidate struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID    string `gorm:"not null;type:uuid"`
	FullName  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobApplication struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
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

type Role struct {
	ID                  string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name                string `gorm:"not null"`
	CreatedAt           time.Time
	OrganizationID      string `gorm:"not null"`
	HomePermission      bool   `gorm:"not null;default:false"`
	CreateJobPermission bool   `gorm:"not null;default:false"`
	ViewJobPermission   bool   `gorm:"not null;default:false"`
	IamPermission       bool   `gorm:"not null;default:false"`
}
