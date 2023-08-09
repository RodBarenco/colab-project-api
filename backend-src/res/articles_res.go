package res

import (
	"time"

	"github.com/google/uuid"
)

type ArticleResponse struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	AuthorName     string      `json:"author_name"`
	Subject        string      `json:"subject"`
	Field          string      `json:"field"`
	Description    string      `json:"description"`
	Keywords       string      `json:"keywords"`
	SubmissionDate time.Time   `json:"submission_date"`
	LikedBy        []uuid.UUID `json:"liked_by"`
	Shares         int         `json:"shares"`
	CoverImage     string      `json:"cover_image"`
}
