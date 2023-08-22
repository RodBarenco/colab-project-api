package res

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
