package dto

type UpdateCommentDto struct {
	Content string `json:"content" validate:"required,min=1"`
}
