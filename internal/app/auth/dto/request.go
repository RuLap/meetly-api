package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type GoogleAuthRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" json:"state"`
}

type ConfirmEmailRequest struct {
	Token string `json:"token" validate:"required,uuid4"`
}
