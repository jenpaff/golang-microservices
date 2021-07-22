package api

type UserCreationRequest struct {
	UserName    string `json:"name" validate:"required,validRegexInput"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
