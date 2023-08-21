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
	PublicKey       string
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
func AddUserToLikedByFromArticle(db *gorm.DB, articleID uuid.UUID, userID uuid.UUID) error {
	article := Article{ID: articleID}
	user := User{ID: userID}
	return db.Model(&article).Association("LikedBy").Append([]User{user})
}

// Função para remover um usuário da lista de "likedBy" de um artigo
func RemoveUserFromLikedByFromArticle(db *gorm.DB, articleID uuid.UUID, userID uuid.UUID) error {
	article := Article{ID: articleID}
	user := User{ID: userID}
	return db.Model(&article).Association("LikedBy").Delete([]User{user})
}

// Função para obter os usuários que "deram like" a um artigo
func GetLikedByUsers(db *gorm.DB, articleID uuid.UUID) ([]User, error) {
	var article Article
	result := db.Preload("LikedBy").Where("id = ?", articleID).First(&article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Article not found")
		}
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

//--------------------------CITATIONS-------------------------------------//

type CitingArticleRequestParams struct {
	CitingArticleID uuid.UUID
	CitedArticleID  uuid.UUID
	UserID          uuid.UUID
}

func AddCitation(db *gorm.DB, citingArticleID, citedArticleID, userID uuid.UUID) error {
	citingArticle := Article{}
	citedArticle := Article{}

	// Verificar se o usuário é o autor do artigo que está citando
	if err := db.First(&citingArticle, "id = ? AND author_id = ?", citingArticleID, userID).Error; err != nil {
		return errors.New("User is not authorized to add citation to this article - or article doesn't exist")
	}

	// Verificar se o artigo citado existe
	if err := db.First(&citedArticle, "id = ?", citedArticleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Cited article not found")
		}
		return err
	}

	return db.Model(&citingArticle).Association("Citations").Append([]Article{citedArticle})
}

func RemoveCitation(db *gorm.DB, citingArticleID, citedArticleID, userID uuid.UUID) error {
	citingArticle := Article{}
	citedArticle := Article{}

	// Verificar se o usuário é o autor do artigo que está citando
	if err := db.First(&citingArticle, "id = ? AND author_id = ?", citingArticleID, userID).Error; err != nil {
		return errors.New("User is not authorized to remove citation from this article - or article doesn't exist")
	}

	// Remover a citação do artigo citado
	if err := db.First(&citedArticle, "id = ?", citedArticleID).Error; err != nil {
		return errors.New("Cited article not found")
	}

	return db.Model(&citingArticle).Association("Citations").Delete([]Article{citedArticle})
}

func GetCitingArticles(db *gorm.DB, articleID uuid.UUID) ([]Article, error) {
	var article Article
	result := db.Preload("Citations").Where("id = ?", articleID).First(&article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Article not found")
		}
		return nil, result.Error
	}
	return article.Citations, nil
}

// if needed
// ArticleCitedBy retrieves a list of articles that cite the specified article.
func ArticleCitedBy(db *gorm.DB, articleID uuid.UUID) ([]Article, error) {
	var citedByArticles []Article
	result := db.Model(&Article{}).
		Where("article_citations.citation_id = ?", articleID).
		Joins("JOIN article_citations ON articles.id = article_citations.article_id").
		Find(&citedByArticles)
	if result.Error != nil {
		return nil, result.Error
	}
	return citedByArticles, nil
}

//-------------------SHARES-----------------------------------------//

func IncrementArticleShares(db *gorm.DB, articleID uuid.UUID) error {
	var article Article
	result := db.Model(&article).Where("id = ?", articleID).UpdateColumn("shares", gorm.Expr("shares + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// --------------------KEY---------------------------//

type AddPublicKeyRequest struct {
	UserID    uuid.UUID
	PublicKey string
}

// AddPublicKeyToUser adiciona uma chave pública ao usuário pelo ID.
func AddPublicKeyToUser(db *gorm.DB, userID uuid.UUID, publicKey string) error {
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	// Atualizar a chave pública do usuário
	user.PublicKey = publicKey
	result = db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserPublicKey retorna a chave pública do usuário pelo ID.
func GetUserPublicKey(db *gorm.DB, userID uuid.UUID) (string, error) {
	var user User
	result := db.Select("public_key").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.PublicKey, nil
}
