package dto

type UserResponseDto struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ShortBio string `json:"shortBio"`
}

type UserResponse struct {
	ID       string       `json:"ID"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	ShortBio string       `json:"shortBio"`
	Session  SessionBrief `json:"session"`
}

type SessionBrief struct {
	RefreshToken string `json:"refreshToken"`
}

type PostResponse struct {
	ID        string `json:"ID"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Author    Author `json:"author"`
	Views     int    `json:"views"`
}

type Author struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
}

type CommentResponse struct {
	ID        string          `json:"ID"`
	Content   string          `json:"content"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
	Author    Author          `json:"author"`
	Replies   []ReplyResponse `json:"replies"`
}

type ReplyResponse struct {
	ID        string `json:"ID"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Author    Author `json:"author"`
}
