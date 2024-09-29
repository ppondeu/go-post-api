package dto

type BookmarkDto struct {
	PostID string `json:"postID" validate:"required,uuid"`
	UserID string `json:"userID" validate:"required,uuid"`
}
