package res

type ArticleCreatedResponse struct {
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Field   string `json:"field"`
	// Add other fields as needed
}

type CreateArticleRes struct {
	Article ArticleCreatedResponse `json:"article"`
	Message string                 `json:"message"`
}
