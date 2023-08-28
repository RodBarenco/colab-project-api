package db

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	Following       []*User `gorm:"many2many:user_following;"`
	ProfilePhoto    string  // URL for profile photo
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

// FOLLOW

type AddUserToFollowing struct {
	UserID      uuid.UUID
	FollowingID uuid.UUID
}

func FollowUser(db *gorm.DB, followerID, followingID uuid.UUID) error {
	follower := User{}
	following := User{}

	// Find the users by their IDs
	result := db.First(&follower, "id = ?", followerID)
	if result.Error != nil {
		return result.Error
	}

	result = db.First(&following, "id = ?", followingID)
	if result.Error != nil {
		return result.Error
	}

	// Add the following user to the follower's Following list
	return db.Model(&follower).Association("Following").Append(&following)

}

func UnfollowUser(db *gorm.DB, followerID, followingID uuid.UUID) error {
	follower := User{}
	following := User{}

	// Find the users by their IDs
	result := db.First(&follower, "id = ?", followerID)
	if result.Error != nil {
		return result.Error
	}

	result = db.First(&following, "id = ?", followingID)
	if result.Error != nil {
		return result.Error
	}

	return db.Model(&follower).Association("Following").Delete(&following)
}

//-----------------------	UPDATE ---------------------------//

type UserUpdateParams struct {
	FirstName     *string
	LastName      *string
	Nickname      *string
	Email         *string
	Password      *string
	DateOfBirth   *string
	Field         *string
	Biography     *string
	Lcourse       *string
	Ccourse       *string
	OpenToColab   *bool
	PublicKey     *string
	ProfilePhoto  *string
	LastEducation *Institution
	Currently     *Institution
}

func UpdateUserFields(db *gorm.DB, userID uuid.UUID, updateParams UserUpdateParams) ([]string, error) {
	user := User{ID: userID}
	fildsUpdated := []string{}

	if updateParams.FirstName != nil {
		user.FirstName = *updateParams.FirstName
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated first name to: %s", *updateParams.FirstName))
	}
	if updateParams.LastName != nil {
		user.LastName = *updateParams.LastName
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated last name to: %s", *updateParams.LastName))
	}
	if updateParams.Nickname != nil {
		user.Nickname = *updateParams.Nickname
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated nick name to: %s", *updateParams.Nickname))
	}
	if updateParams.Biography != nil {
		user.Biography = *updateParams.Biography
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated Biography to: %s", *updateParams.Biography))
	}
	if updateParams.DateOfBirth != nil {
		dateOfBirth, err := time.Parse("2006-01-02", *updateParams.DateOfBirth)
		if err != nil {
			return nil, err
		}
		user.DateOfBirth = dateOfBirth
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated birth date to: %s", *updateParams.DateOfBirth))
	}
	if updateParams.Email != nil {
		user.Email = *updateParams.Email
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated email to: %s", *updateParams.Email))
	}
	if updateParams.Password != nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateParams.Password), bcrypt.DefaultCost)
		if err != nil {
			// Return a more specific error message including the original password
			return nil, err
		}
		user.Password = string(hashedPassword)
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated password to: %s", *updateParams.Password))
	}
	if updateParams.Field != nil {
		user.Field = *updateParams.Field
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated field to: %s", *updateParams.Field))
	}
	if updateParams.Ccourse != nil {
		user.Ccourse = *updateParams.Ccourse
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated current course to: %s", *updateParams.Ccourse))
	}
	if updateParams.Lcourse != nil {
		user.Lcourse = *updateParams.Lcourse
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated last course to: %s", *updateParams.Lcourse))
	}
	if updateParams.OpenToColab != nil {
		user.OpenToColab = *updateParams.OpenToColab
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated open to colab to: %v", strconv.FormatBool(*updateParams.OpenToColab)))
	}
	if updateParams.ProfilePhoto != nil {
		user.ProfilePhoto = *updateParams.ProfilePhoto
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated porfile photo to: %s", *updateParams.ProfilePhoto))
	}
	if updateParams.PublicKey != nil {
		user.PublicKey = *updateParams.PublicKey
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Update publick keyto: %s", *updateParams.PublicKey))
	}
	if updateParams.LastEducation != nil {
		user.LastEducation = *updateParams.LastEducation
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated last educatinal institution to: %s", *&updateParams.LastEducation.Name))
	}
	if updateParams.Currently != nil {
		user.Currently = *updateParams.Currently
		fildsUpdated = append(fildsUpdated, fmt.Sprintf("Updated currently educational institution to: %s", *&updateParams.Currently.Name))
	}

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return fildsUpdated, nil
}
