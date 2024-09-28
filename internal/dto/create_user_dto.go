package dto

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,lowercase"`
	Password string `json:"password" validate:"required,min=6"`
	ShortBio string `json:"short_bio" validate:"max=160"`
}
