package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func GetUserByID(db *gorm.DB, userID uuid.UUID) (User, error) {
	var user User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// -------------------------- LIKES ----------------------------//

type LikeArticleRequestParams struct {
	UserID    uuid.UUID
	ArticleID uuid.UUID
}

// Função para adicionar um usuário que "deu like" a um artigo
func AddUserToLikedBy(db *gorm.DB, articleID uuid.UUID, userID uuid.UUID) error {
	article := Article{ID: articleID}
	user := User{ID: userID}
	return db.Model(&article).Association("LikedBy").Append([]User{user})
}

// Função para remover um usuário da lista de "likedBy" de um artigo
func RemoveUserFromLikedBy(db *gorm.DB, articleID uuid.UUID, userID uuid.UUID) error {
	article := Article{ID: articleID}
	user := User{ID: userID}
	return db.Model(&article).Association("LikedBy").Delete([]User{user})
}

// Função para obter os usuários que "deram like" a um artigo
func GetLikedByUsers(db *gorm.DB, articleID uuid.UUID) ([]User, error) {
	var article Article
	result := db.Preload("LikedBy").Where("id = ?", articleID).First(&article)
	if result.Error != nil {
		return nil, result.Error
	}
	return article.LikedBy, nil
}

func IsArticleLikedByUser(dbAccess *gorm.DB, articleID uuid.UUID, userID uuid.UUID) (bool, error) {
	var article Article
	result := dbAccess.Preload("LikedBy").
		Where("id = ?", articleID).
		First(&article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	for _, likedByUser := range article.LikedBy {
		if likedByUser.ID == userID {
			return true, nil
		}
	}

	return false, nil
}
