package repository

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/domain"
	"gorm.io/gorm"
)

type FollowRepositoryDB struct {
	db *gorm.DB
}

func NewFollowRepositoryDB(db *gorm.DB) domain.FollowRepository {
	return &FollowRepositoryDB{db}
}

func (r *FollowRepositoryDB) Create(Follow *domain.Follow) (*domain.Follow, error) {
	if err := r.db.Create(Follow).Error; err != nil {
		return nil, err
	}
	return Follow, nil
}

func (r *FollowRepositoryDB) FindAll() ([]domain.Follow, error) {
	var Follows []domain.Follow
	if err := r.db.Find(&Follows).Error; err != nil {
		return nil, err
	}
	return Follows, nil
}

func (r *FollowRepositoryDB) FindByID(ID uuid.UUID) (*domain.Follow, error) {
	var Follow domain.Follow
	if err := r.db.Where("id = ?", ID).First(&Follow).Error; err != nil {
		return nil, err
	}
	return &Follow, nil
}

func (r *FollowRepositoryDB) FindByFollowerIDAndFollowedID(FollowerID, followedID uuid.UUID) (*domain.Follow, error) {
	var Follow domain.Follow
	if err := r.db.Where("follower_id = ? AND followed_id = ?", FollowerID, followedID).First(&Follow).Error; err != nil {
		return nil, err
	}
	return &Follow, nil
}

func (r *FollowRepositoryDB) Delete(ID uuid.UUID) error {
	if err := r.db.Where("id = ?", ID).Delete(&domain.Follow{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *FollowRepositoryDB) FindFollowersByUserID(userID uuid.UUID) ([]domain.Follow, error) {
	var Follows []domain.Follow
	if err := r.db.Where("followed_id = ?", userID).Find(&Follows).Error; err != nil {
		return nil, err
	}
	return Follows, nil
}

func (r *FollowRepositoryDB) FindFollowedUsersByUserID(userID uuid.UUID) ([]domain.Follow, error) {
	var Follows []domain.Follow
	if err := r.db.Where("follower_id = ?", userID).Find(&Follows).Error; err != nil {
		return nil, err
	}
	return Follows, nil
}
