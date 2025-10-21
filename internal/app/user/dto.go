package user

type GetUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"Gender"`
	AvatarUrl string `json:"avatar_url"`
}

type SaveUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	BirthDate string `json:"birth_date" validate:"required,date"`
	Gender    string `json:"gender" validate:"required"`
	AvatarUrl string `json:"avatar_url" validate:"required,url"`
}
