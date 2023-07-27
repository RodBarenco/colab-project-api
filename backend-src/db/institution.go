package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Institution struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"not null"`
}

func IsInstitutionExists(db *gorm.DB, id uuid.UUID) bool {
	var institution Institution
	result := db.Where("id = ?", id).First(&institution)
	return result.RowsAffected > 0 && result.Error == nil
}
