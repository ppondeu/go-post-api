package repository

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/domain"
	"github.com/ppondeu/go-post-api/internal/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepositoryDB struct {
	db *gorm.DB
}

func NewPostRepositoryDB(db *gorm.DB) domain.PostRepository {
	return &PostRepositoryDB{db}
}

func (r *PostRepositoryDB) FindAll() ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Preload(clause.Associations).Find(&posts).Error
	if err != nil {
		return nil, err
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
	err = r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Preload("Likes").Preload("Comments").First(&savedPost, ID).Error
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
	err = r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Preload("Likes").Preload("Comments").First(&updatedPost, ID).Error

	if err != nil {
		return nil, err
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

func (r *PostRepositoryDB) LikePost(like domain.Like) error {
	err := r.db.Create(&like).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepositoryDB) UnlikePost(userID, postID uuid.UUID) error {
	err := r.db.Delete(&domain.Like{}, "user_id = ? AND post_id = ?", userID, postID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepositoryDB) GetPostLikeCount(postID uuid.UUID) (uint32, error) {
	var count uint32
	int64count := int64(count)
	result := r.db.Model(&domain.Like{}).Where("post_id = ?", postID).Count(&int64count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *PostRepositoryDB) AddComment(comment domain.Comment) (*domain.Comment, error) {
	err := r.db.Create(&comment).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var savedComment domain.Comment
	ID, _ := uuid.Parse(comment.ID)
	err = r.db.Preload("Replies", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, content, user_id")
	}).First(&savedComment, ID).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &savedComment, nil
}

func (r *PostRepositoryDB) UpdateComment(ID uuid.UUID, comment domain.Comment) (*domain.Comment, error) {
	err := r.db.Model(&domain.Comment{}).Where("id = ?", ID).Updates(comment).Error
	if err != nil {
		return nil, err
	}

	var updatedComment domain.Comment
	err = r.db.Preload("Replies", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, content, user_id, post_id, parent_id")
	}).First(&updatedComment, ID).Error
	if err != nil {
		return nil, err
	}

	return &updatedComment, nil
}

func (r *PostRepositoryDB) DeleteComment(ID uuid.UUID) error {
	result := r.db.Delete(&domain.Comment{}, ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostRepositoryDB) FindCommentsByPostID(postID uuid.UUID) ([]domain.Comment, error) {
	var comments []domain.Comment
	result := r.db.Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (r *PostRepositoryDB) FindCommentByID(ID uuid.UUID) (*domain.Comment, error) {
	var comment domain.Comment
	result := r.db.Preload("Replies", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, content, user_id, post_id, parent_id")
	}).First(&comment, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (r *PostRepositoryDB) FindRepliesByCommentID(commentID uuid.UUID) ([]domain.Comment, error) {
	var comments []domain.Comment
	result := r.db.Select("id, content, user_id, post_id, parent_id").Where("parent_id = ?", commentID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
