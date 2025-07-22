package models

import (
	"time"

	"github.com/lib/pq"
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
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID     string `gorm:"not null;type:uuid"`
	FullName   string `gorm:"not null"`
	Email      string `gorm:"not null"`
	Phone      string `gorm:"not null"`
	LinkedIn   string `gorm:"type:text"`
	GitHub     string `gorm:"type:text"`
	Location   string `gorm:"type:text"`
	Experience string `gorm:"type:text"`
	Education  string `gorm:"type:text"`
	Skills     string `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type JobApplication struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CandidateID string `gorm:"type:uuid"`
	JobID       string `gorm:"type:uuid"`

	ResumeGCSPath      string  `gorm:"not null"`
	CoverLetterGCSPath string  `gorm:"type:text"`
	ParsedResume       string  `gorm:"type:jsonb"`
	PinecodeID         string  `gorm:"type:text"` // pinecone id for embedding retrieval
	Status             string  `gorm:"not null;default:'pending'"`
	AI_Score           float64 `gorm:"not null;default:0"`
	MagicLinkToken     string  `gorm:"unique;not null"`

	CreatedAt time.Time
}

type Job struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrganizationID string `gorm:"type:uuid;not null"`

	Title           string         `gorm:"not null"`
	Description     string         `gorm:"not null;type:text"`
	Location        pq.StringArray `gorm:"type:text[]"`
	ExperienceLevel string         `gorm:"type:text"`
	SkillsRequired  pq.StringArray `gorm:"type:text[]"`
	EmploymentType  pq.StringArray `gorm:"type:text[]"`
	SalaryRange     pq.StringArray `gorm:"type:text[]"`
	IsActive        bool           `gorm:"not null;default:true"`

	CreatedByID      string `gorm:"type:uuid;not null"`
	ApplicationCount int    `gorm:"default:0"`

	PublicLink string `gorm:"unique;not null"` // resumelens.com/job/{org_id}/{job_id}
	ShortLink  string `gorm:"unique"`          // resumelens.com/job/{job_id}
	CreatedAt  time.Time
}

type JobAnalytics struct {
	ID    string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	JobID string `gorm:"type:uuid"`

	TotalApplications int     `gorm:"not null;default:0"`
	TotalHires        int     `gorm:"not null;default:0"`
	AvgFitScore       float64 `gorm:"not null;default:0"`
	CreatedAt         time.Time
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
