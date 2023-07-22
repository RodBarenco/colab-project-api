package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName     string    `gorm:"not null"`
	LastName      string    `gorm:"not null"`
	Nickname      string
	Email         string    `gorm:"not null;unique"`
	DateOfBirth   time.Time `gorm:"not null"`
	Field         string
	Interests     string
	Biography     string
	LastEducation string
	Currently     string
	Institution   Institution `gorm:"foreignKey:InstitutionID"`
	InstitutionID uuid.UUID
	OpenToColab   bool
}
