package res

import (
	"time"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/google/uuid"
)

type ArticleResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	AuthorName     string    `json:"author_name"`
	Subject        string    `json:"subject"`
	Field          string    `json:"field"`
	Description    string    `json:"description"`
	Keywords       string    `json:"keywords"`
	SubmissionDate time.Time `json:"submission_date"`
	LikedBy        []string  `json:"liked_by"`
	Shares         int       `json:"shares"`
	CoverImage     string    `json:"cover_image"`
}

// -------------------------------------------------------
type LikeArticleResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	ArticleID uuid.UUID `json:"artcle_id"`
}

type CreateLikeArticleResponse struct {
	LikeArticleResponse `json:"article_liked_by"`
	Message             string `json:"message"`
}

// ----------------------------------------------------------
type ArticleWithLikesResponse struct {
	Article       db.Article `json:"article"`
	RelatedTables LikesInfo  `json:"relatedTables"`
	Message       string     `json:"message"`
}

type LikesInfo struct {
	NumLikes     int      `json:"numLikes"`
	LikedByNames []string `json:"likedByNames"`
}

// ----------------------------------------------------------

type LikedByUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type LikedByUsersResponse struct {
	LikedByUsers []LikedByUser `json:"liked_by_users"`
	Message      string        `json:"message"`
}

// -------------------------------------
// ArticleCitingCitedRes defines the structure for the response containing article citing/cited information.
type ArticleCitingCitedRes struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
}
