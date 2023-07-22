package db

import (
	"github.com/google/uuid"
)

type Institution struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"not null"`
}
