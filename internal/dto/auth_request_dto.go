package dto

type AuthRequestDTO struct {
	Email    string `json:"email" validate:"required,email,min=4,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}
