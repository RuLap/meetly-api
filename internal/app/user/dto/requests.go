package dto

type SaveUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	BirthDate string `json:"birth_date" validate:"required,date"`
	Gender    string `json:"gender" validate:"required"`
	AvatarUrl string `json:"avatar_url" validate:"required,url"`
}
