package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	IsAccepted     bool
	LikedBy        []uuid.UUID `gorm:"type:uuid[]"`
	Citations      []uuid.UUID `gorm:"type:uuid[]"`
	Shares         int
	CoAuthors      string
	CoverImage     string
}

type ArticleSearchParams struct {
	Search []string `json:"search"`
}

type ArticleParams struct {
	Title          string
	AuthorID       uuid.UUID
	Subject        string
	Field          string
	File           []byte
	Description    string
	Keywords       string
	SubmissionDate time.Time
	LikedBy        []uuid.UUID
	Citations      []uuid.UUID
	Shares         int
	CoAuthors      string
	CoverImage     string
}

func CreateArticle(db *gorm.DB, newArticle ArticleParams) error {
	newArticle.SubmissionDate = time.Now()
	article := Article{
		Title:          newArticle.Title,
		AuthorID:       newArticle.AuthorID,
		Subject:        newArticle.Subject,
		Field:          newArticle.Field,
		File:           newArticle.File,
		Description:    newArticle.Description,
		Keywords:       newArticle.Keywords,
		SubmissionDate: newArticle.SubmissionDate,
		CoAuthors:      newArticle.CoAuthors,
		CoverImage:     newArticle.CoverImage,
	}

	return db.Create(&article).Error
}

func GetLatestThousandArticles(db *gorm.DB) ([]Article, error) {
	var articles []Article
	err := db.Order("submission_date desc").Limit(1000).Find(&articles).Error
	return articles, err
}

// Obter os 50 artigos mais recentes
func GetLatestFiftyArticles(db *gorm.DB) ([]Article, error) {
	var articles []Article
	err := db.Order("submission_date desc").Limit(50).Find(&articles).Error
	return articles, err
}

// Obter artigos filtrados por título
func GetArticlesByTitle(db *gorm.DB, title string) ([]Article, error) {
	var articles []Article
	err := db.Where("title = ?", title).Find(&articles).Error
	return articles, err
}

// Obter artigos filtrados por tema
func GetArticlesBySubject(db *gorm.DB, subject string) ([]Article, error) {
	var articles []Article
	err := db.Where("subject = ?", subject).Find(&articles).Error
	return articles, err
}

// Obter artigos filtrados por autor
func GetArticlesByAuthor(db *gorm.DB, authorID uuid.UUID) ([]Article, error) {
	var articles []Article
	err := db.Where("author_id = ?", authorID).Find(&articles).Error
	return articles, err
}

// Obter artigos filtrados por campo
func GetArticlesByField(db *gorm.DB, field string) ([]Article, error) {
	var articles []Article
	err := db.Where("field = ?", field).Find(&articles).Error
	return articles, err
}

// Obter artigos filtrados por palavras-chave
func GetArticlesByKeywords(db *gorm.DB, keywords ...string) ([]Article, error) {
	var articles []Article
	for _, keyword := range keywords {
		var filteredArticles []Article
		err := db.Where("keywords LIKE ?", "%"+keyword+"%").Find(&filteredArticles).Error
		if err != nil {
			return nil, err
		}
		articles = append(articles, filteredArticles...)
	}
	return articles, nil
}

// Obter o artigo por ID
func GetLatestArticleById(db *gorm.DB, articleID uuid.UUID) (Article, error) {
	var article Article
	err := db.Where("id = ?", articleID).Order("submission_date desc").First(&article).Error
	return article, err
}
