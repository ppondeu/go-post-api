package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Comment struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Content   string    `gorm:"type:varchar(255);not null" json:"content"`
	UserID    string    `gorm:"type:uuid;not null;" json:"user_id"`
	PostID    string    `gorm:"type:uuid;not null;" json:"post_id"`
	ParentID  *string   `gorm:"type:uuid;index" json:"parent_id"`
	Replies   []Comment `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"replies"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp;autoUpdateTime" json:"updated_at"`
}

type Like struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_like" json:"post_id"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

type Bookmark struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_bookmark" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

type Tag struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string `gorm:"type:varchar(255);not null;unique"`
}

type Post struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   string         `gorm:"not null" json:"content"`
	Views     int            `gorm:"default:0" json:"views"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	UserID    string         `gorm:"type:uuid;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Likes     []Like         `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"likes"`
	Bookmarks []Bookmark     `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bookmarks"`
	Comments  []Comment      `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
}

type PostRepository interface {
	FindAll() ([]Post, error)
	FindByID(ID uuid.UUID) (*Post, error)
	FindByUserID(userID uuid.UUID) ([]Post, error)
	Save(post Post) (*Post, error)
	Update(ID uuid.UUID, post Post) (*Post, error)
	Delete(ID uuid.UUID) error

	FindAllTags() ([]Tag, error)
	AddBookmark(bookmark Bookmark) error
	RemoveBookmark(userID, postID uuid.UUID) error

	CreateTag(tag Tag) (*Tag, error)

	LikePost(like Like) error
	UnlikePost(userID, postID uuid.UUID) error
}
