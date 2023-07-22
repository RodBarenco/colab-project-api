package db

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title          string    `gorm:"not null"`
	AuthorID       uuid.UUID `gorm:"not null;foreignKey:AuthorID"`
	Subject        string    `gorm:"not null"`
	Field          string    `gorm:"not null"`
	File           []byte    `gorm:"not null"`
	Description    string    `gorm:"not null"`
	Keywords       string    `gorm:"not null"`
	SubmissionDate time.Time `gorm:"not null"`
	AcceptanceDate time.Time
	LikedBy        []uuid.UUID `gorm:"type:uuid[]"`
	Citations      []uuid.UUID `gorm:"type:uuid[]"`
	Shares         int
	CoAuthors      string
	CoverImage     string
}
