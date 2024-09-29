package dto

type UpdatePostDto struct {
	Title   string   `json:"title" validate:"omitempty"`
	Content string   `json:"content" validate:"omitempty,min=3"`
	Tags    []string `json:"tags" validate:"omitempty,dive,required,min=1"`
}
