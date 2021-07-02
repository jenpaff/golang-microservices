package api

type UserCreationRequest struct {
	UserName    string `json:"name" validate:"notBlank,validRegexInput"`
	Email       string `json:"email" validate:"notBlank"`
	PhoneNumber string `json:"phone_number" validate:"notBlank"`
}
