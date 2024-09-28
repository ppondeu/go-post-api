package domain

import (
	"time"

	"github.com/lib/pq"
)

type Comment struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Content         string    `gorm:"type:varchar(255);not null"`
	UserID          string    `gorm:"type:uuid;not null"`
	PostID          string    `gorm:"type:uuid;not null"`
	ParentCommentID *string   `gorm:"type:uuid"`
	Replies         []Comment `gorm:"foreignKey:ParentCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp;autoUpdateTime"`
}

type Like struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_like"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_like"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

type Bookmark struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_bookmark"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_bookmark"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

type Tag struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string `gorm:"type:varchar(255);not null"`
}

type Post struct {
	ID       string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title    string         `gorm:"type:varchar(255);not null"`
	Content  string         `gorm:"not null"`
	Views    int            `gorm:"default:0"`
	Tags     pq.StringArray `gorm:"type:text[]"`
	UserID   string         `gorm:"type:uuid;not null"`
	User     User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Likes    []Like         `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Comments []Comment      `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
