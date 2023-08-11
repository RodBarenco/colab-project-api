package utils

import (
	"github.com/RodBarenco/colab-project-api/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAuthorName(dbAccess *gorm.DB, authorID uuid.UUID) (string, error) {
	var user db.User
	err := dbAccess.Where("id = ?", authorID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.FirstName + " " + user.LastName, nil
}
