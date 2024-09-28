package dto

type UpdateUserDto struct {
	Username string `json:"username" validate:"omitempty,min=3,lowercase"`
	Password string `json:"password" validate:"omitempty,min=6"`
	ShortBio string `json:"short_bio" validate:"omitempty,max=160"`
}
