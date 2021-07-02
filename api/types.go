package api

type UserCreationRequest struct {
	UserName    string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
