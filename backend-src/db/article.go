package db

import (
	"strings"
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
	LikedBy        []User    `gorm:"many2many:article_likes"`
	Citations      []Article `gorm:"many2many:article_citations"`
	Shares         int
	CoAuthors      string
	CoverImage     string
	ApprovedBy     []*Admin `gorm:"many2many:approved_by;"`
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

	newArticle.Field = strings.ToLower(newArticle.Field)
	newArticle.Subject = strings.ToLower(newArticle.Subject)
	newArticle.Keywords = strings.ToLower(newArticle.Keywords)

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
	subject = strings.ToLower(subject)
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
	field = strings.ToLower(field)
	err := db.Where("field = ?", field).Find(&articles).Error
	return articles, err
}

// Obter artigos filtrados por palavras-chave
func GetArticlesByKeywords(db *gorm.DB, keywords ...string) ([]Article, error) {
	var articles []Article
	for _, keyword := range keywords {
		keyword = strings.ToLower(keyword)
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
func GetArticleById(db *gorm.DB, articleID uuid.UUID) (Article, error) {
	var article Article
	err := db.Where("id = ?", articleID).First(&article).Error
	return article, err
}

// Articles Recomended for a specificc user

func GetRecommendedArticles(db *gorm.DB, userID uuid.UUID, monthsAgo int) ([]Article, []Article, error) {
	var user User
	if err := db.Preload("Interests").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, nil, err
	}

	var articles []Article
	var otherArticles []Article

	var interestNames []string
	for _, interest := range user.Interests {
		interestNames = append(interestNames, interest.Name)
	}

	// Calculate the date 'monthsAgo' months ago from now
	var cutoffDate time.Time
	if monthsAgo > 0 {
		cutoffDate = time.Now().AddDate(0, -monthsAgo, 0)
	}

	// Select articles with matching interests or matching field, order by likes
	query := db.
		Model(&Article{}).
		Where("is_accepted = ?", true).
		Where("(field ILIKE ? OR exists(select 1 from interests, user_interests where user_interests.interest_id = interests.id and user_interests.user_id = ? and interests.name ILIKE any(array[?])))", "%"+user.Field+"%", userID, interestNames).
		Order("array_length(liked_by, 1) DESC").
		Limit(50)

	if !cutoffDate.IsZero() {
		query = query.Where("submission_date >= ?", cutoffDate)
	}

	err := query.Find(&articles).Error
	if err != nil {
		return nil, nil, err
	}

	// Create a map to keep track of article IDs
	articleIDMap := make(map[uuid.UUID]bool)

	// Iterate through the recommended articles and add their IDs to the map
	for _, article := range articles {
		articleIDMap[article.ID] = true
	}

	// If not enough articles found, fetch latest articles to complete the list
	limit := 50 - len(articles)
	if limit > 0 {
		otherQuery := db.
			Model(&Article{}).
			Where("is_accepted = ?", true).
			Order("submission_date DESC").
			Limit(limit)

		if !cutoffDate.IsZero() {
			otherQuery = otherQuery.Where("submission_date >= ?", cutoffDate)
		}

		// Restrict articles to a specific date range
		if !cutoffDate.IsZero() {
			otherQuery = otherQuery.Where("submission_date >= ?", cutoffDate)
		}

		err := otherQuery.Find(&otherArticles).Error
		if err != nil {
			return nil, nil, err
		}

		// Iterate through the other articles and add to otherArticleResponses
		var filteredOtherArticles []Article
		for _, article := range otherArticles {
			// Check if the article ID is already in the map
			if _, ok := articleIDMap[article.ID]; !ok {
				filteredOtherArticles = append(filteredOtherArticles, article)
			}
		}
		otherArticles = filteredOtherArticles
	}

	return articles, otherArticles, nil
}
