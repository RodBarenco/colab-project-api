package utils

import (
	"github.com/RodBarenco/colab-project-api/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func GetAuthorIDByName(dbAccess *gorm.DB, firstName, lastName string) (uuid.UUID, error) {
	var user db.User
	err := dbAccess.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&user).Error
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func GetNamesOfUsersThatLikedArticles(dbAcess *gorm.DB, articleID uuid.UUID) ([]db.User, error) {
	var article db.Article
	result := dbAcess.Preload("LikedBy").Where("id = ?", articleID).First(&article)
	if result.Error != nil {
		return nil, result.Error
	}
	return article.LikedBy, nil
}

func GetUserIDByEmail(dbAccessor *gorm.DB, email string) (uuid.UUID, error) {
	var user db.User
	if err := dbAccessor.Where("email = ?", email).First(&user).Error; err != nil {
		return uuid.UUID{}, err
	}
	return user.ID, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GetAdminIDByEmail(dbAccessor *gorm.DB, email string) (uuid.UUID, error) {
	var admin db.Admin
	if err := dbAccessor.Where("email = ?", email).First(&admin).Error; err != nil {
		return uuid.UUID{}, err
	}
	return admin.ID, nil
}
