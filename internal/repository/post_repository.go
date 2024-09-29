package repository

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/domain"
	"gorm.io/gorm"
)

type PostRepositoryDB struct {
	db *gorm.DB
}

func NewPostRepositoryDB(db *gorm.DB) domain.PostRepository {
	return &PostRepositoryDB{db}
}

func (r *PostRepositoryDB) FindAll() ([]domain.Post, error) {
	var posts []domain.Post
	result := r.db.Preload("User").Preload("Likes").Preload("Comments").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (r *PostRepositoryDB) FindByID(ID uuid.UUID) (*domain.Post, error) {
	var post domain.Post
	result := r.db.Preload("User").Preload("Likes").Preload("Comments").First(&post, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (r *PostRepositoryDB) FindByUserID(userID uuid.UUID) ([]domain.Post, error) {
	var posts []domain.Post
	result := r.db.Preload("User").Preload("Likes").Preload("Comments").Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (r *PostRepositoryDB) Save(post domain.Post) (*domain.Post, error) {
	err := r.db.Create(&post).Error
	if err != nil {
		return nil, err
	}
	var savedPost domain.Post
	ID, _ := uuid.Parse(post.ID)
	err = r.db.Preload("User").Preload("Likes").Preload("Comments").First(&savedPost, ID).Error
	if err != nil {
		return nil, err
	}
	return &savedPost, nil
}

func (r *PostRepositoryDB) Update(ID uuid.UUID, post domain.Post) (*domain.Post, error) {
	err := r.db.Model(&domain.Post{}).Where("id = ?", ID).Updates(post).Error
	if err != nil {
		return nil, err
	}

	var updatedPost domain.Post
	result := r.db.Preload("User").Preload("Likes").Preload("Comments").First(&updatedPost, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedPost, nil
}

func (r *PostRepositoryDB) Delete(ID uuid.UUID) error {
	result := r.db.Delete(&domain.Post{}, ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostRepositoryDB) FindAllTags() ([]domain.Tag, error) {
	var tags []domain.Tag
	result := r.db.Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (r *PostRepositoryDB) AddBookmark(bookmark domain.Bookmark) error {
	result := r.db.Create(&bookmark)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostRepositoryDB) RemoveBookmark(userID, postID uuid.UUID) error {
	result := r.db.Delete(&domain.Bookmark{}, "user_id = ? AND post_id = ?", userID, postID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostRepositoryDB) CreateTag(tag domain.Tag) (*domain.Tag, error) {
	err := r.db.Create(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
