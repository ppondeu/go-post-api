package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username    string      `gorm:"unique;not null" json:"username"`
	Email       string      `gorm:"unique;not null" json:"email"`
	Password    string      `gorm:"not null" json:"password"`
	ShortBio    string      `gorm:"type:varchar(160);default:''" json:"short_bio"`
	UserSession UserSession `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_session"`
	Posts       []Post      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"posts"`
	Follower    []Follow    `gorm:"foreignKey:FollowerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"follower"`
	Followed    []Follow    `gorm:"foreignKey:FollowedID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"followed"`
}

type UserSession struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       string    `gorm:"type:uuid;not null;unique" json:"user_id"`
	RefreshToken *string   `gorm:"type:text;unique" json:"refresh_token"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp;autoUpdateTime" json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindAll() ([]User, error)
	FindByID(ID uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindUserWithRelation(ID uuid.UUID) (*User, error)
	FindAllUsersWithRelation() ([]User, error)
	Update(ID uuid.UUID, user *User) (*User, error)
	Delete(ID uuid.UUID) error
	CreateUserAndSession(user *User, refreshToken *string) (*User, error)
	UpdateSession(userID uuid.UUID, refreshToken *string) error
	FindSession(userID uuid.UUID) (*UserSession, error)
}
