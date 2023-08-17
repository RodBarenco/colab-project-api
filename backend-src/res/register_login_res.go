package res

import "github.com/google/uuid"

// UserGetedResponse contains the sensitive user information.
type UserGetedResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	// Add other fields as needed
}

type SignupRes struct {
	User    UserGetedResponse `json:"user"`
	Message string            `json:"message"`
}

type LoginRes struct {
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	PublicKey string    `json:"public_key"`
	Message   string    `json:"message"`
}
