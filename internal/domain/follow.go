package domain

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"ID"`
	FollowerID string    `gorm:"type:uuid;not null" json:"followerID"`
	FollowedID string    `gorm:"type:uuid;not null" json:"follwedID"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}

type FollowRepository interface {
	Create(Follow *Follow) (*Follow, error)
	FindAll() ([]Follow, error)
	FindByID(ID uuid.UUID) (*Follow, error)
	FindByFollowerIDAndFollowedID(followerID, followedID uuid.UUID) (*Follow, error)
	Delete(ID uuid.UUID) error
	FindFollowersByUserID(userID uuid.UUID) ([]Follow, error)
	FindFollowedUsersByUserID(userID uuid.UUID) ([]Follow, error)
}
