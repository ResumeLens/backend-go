package models

import "time"

type User struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string `gorm:"unique;not null"`
	PasswordHash   string `gorm:"not null"`
	RoleID         string `gorm:"not null"`
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

type Role struct {
	ID                  string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name                string `gorm:"not null"`
	CreatedAt           time.Time
	OrganizationID      string `gorm:"not null"`
	HomePermission      bool
	CreateJobPermission bool
	ViewJobPermission   bool
	IAMPermission       bool
}

type Invite struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string `gorm:"not null"`
	OrganizationID string `gorm:"not null"`
	RoleID         string `gorm:"not null"`
	Token          string `gorm:"unique;not null"`
	Expiry         time.Time
	IsAccepted     bool
	CreatedAt      time.Time
}
