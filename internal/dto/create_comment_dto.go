package dto

type CreateCommentDto struct {
	Content  string  `json:"content" validate:"required,min=1"`
	UserID   string  `json:"userID" validate:"required,uuid"`
	PostID   string  `json:"postID" validate:"required,uuid"`
	ParentID *string `json:"parentID" validate:"omitempty,uuid"`
}
