package dto

type CreatePostDto struct {
	UserID  string   `json:"userID" validate:"required,uuid4"`
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required,min=3"`
	Tags    []string `json:"tags" validate:"omitempty,dive,required,min=1"`
}
