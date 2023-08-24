package res

import (
	"time"
)

type ArticleCreatedResponse struct {
	Title      string `json:"title"`
	Subject    string `json:"subject"`
	Field      string `json:"field"`
	CoverImage string `json:"conver_image"`
	// Add other fields as needed
}

type CreateArticleRes struct {
	Article ArticleCreatedResponse `json:"article"`
	Message string                 `json:"message"`
}

//--------------------------CITING-----------

// ArticleCitationResponse is the standard response structure for citation actions.
type ArticleCitationResponse struct {
	CitedArticleID  string `json:"cited_article_id"`
	CitingArticleID string `json:"citing_article_id"`
	Message         string `json:"message"`
}

// ------------------------ FOLLOW -----------------------

type FollowUserResponse struct {
	UserID      string `json:"user_id"`
	FollowingID string `json:"following_id"`
	Message     string `json:"message"`
}

//----------------GET USER -------------------

type GetUserResponse struct {
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Nickname      string    `json:"nickname"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	Field         string    `json:"field"`
	Interests     []string  `json:"interests"`
	Biography     string    `json:"biography"`
	LastEducation string    `json:"lastEducation"`
	Lcourse       string    `json:"lcourse"`
	Currently     string    `json:"currently"`
	Ccourse       string    `json:"ccourse"`
	OpenToColab   bool      `json:"openToColab"`
	Following     []string  `json:"following"`
	ProfilePhoto  string    `json:"profilePhoto"`
}
