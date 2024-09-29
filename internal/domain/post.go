package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Comment struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Content   string    `gorm:"type:varchar(255);not null" json:"content"`
	UserID    string    `gorm:"type:uuid;not null;" json:"userID"`
	PostID    string    `gorm:"type:uuid;not null;" json:"postID"`
	ParentID  *string   `gorm:"type:uuid;index" json:"parentID"`
	Replies   []Comment `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"replies"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp;autoUpdateTime" json:"updatedAt"`
}

type Like struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null" json:"userID"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_like" json:"postID"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}

type Bookmark struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null" json:"userID"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_bookmark" json:"postID"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}

type Tag struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string `gorm:"type:varchar(255);not null;unique" json:"name"`
}

type Post struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   string         `gorm:"not null" json:"content"`
	ViewCount int            `gorm:"default:0" json:"viewCount"`
	Tags      pq.StringArray `gorm:"type:text[];default:'{}'" json:"tags"`
	UserID    string         `gorm:"type:uuid;not null" json:"userID"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Likes     []Like         `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"likes,omitempty"`
	Bookmarks []Bookmark     `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bookmarks,omitempty"`
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
	GetPostLikeCount(postID uuid.UUID) (uint32, error)

	AddComment(comment Comment) (*Comment, error)
	UpdateComment(ID uuid.UUID, comment Comment) (*Comment, error)
	DeleteComment(ID uuid.UUID) error
	FindCommentsByPostID(postID uuid.UUID) ([]Comment, error)
	FindCommentByID(ID uuid.UUID) (*Comment, error)
	FindRepliesByCommentID(commentID uuid.UUID) ([]Comment, error)
}
