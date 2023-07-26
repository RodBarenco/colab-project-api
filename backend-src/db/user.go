package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName       string    `gorm:"not null"`
	LastName        string    `gorm:"not null"`
	Nickname        string
	Email           string    `gorm:"not null;unique"`
	Password        string    `gorm:"not null"`
	DateOfBirth     time.Time `gorm:"not null"`
	Field           string
	Interests       []*Interest `gorm:"many2many:user_interests;"`
	Biography       string      `gorm:"type:TEXT"`
	LastEducation   Institution `gorm:"foreignKey:LastEducationID"`
	Lcourse         string
	Currently       Institution `gorm:"foreignKey:CurrentlyID"`
	Ccourse         string
	LastEducationID *uuid.UUID `gorm:"null;type:uuid"`
	CurrentlyID     *uuid.UUID `gorm:"null;type:uuid"`
	OpenToColab     bool
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
