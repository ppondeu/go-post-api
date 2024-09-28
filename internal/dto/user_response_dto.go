package dto

type UserResponseDto struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ShortBio string `json:"short_bio"`
}
