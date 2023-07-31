package res

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
	Message string `json:"message"`
	Token   string `json:"token"`
}
