package repository

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) domain.UserRepository {
	return &UserRepositoryDB{db}
}

func (r *UserRepositoryDB) Create(user *domain.User) (*domain.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryDB) FindAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Preload("UserSession").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryDB) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("UserSession").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("UserSession").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) FindByID(ID uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("UserSession").Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) Update(ID uuid.UUID, user *domain.User) (*domain.User, error) {
	err := r.db.Model(&domain.User{}).Where("id = ?", ID).Updates(user).Error
	if err != nil {
		return nil, err
	}

	var updatedUser domain.User
	err = r.db.Preload("UserSession").Where("id = ?", ID).First(&updatedUser).Error
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *UserRepositoryDB) Delete(ID uuid.UUID) error {
	err := r.db.Where("id = ?", ID).Delete(&domain.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryDB) UpdateSession(userID uuid.UUID, refreshToken *string) error {
	err := r.db.Model(&domain.UserSession{}).
		Preload("UserSession").
		Where("user_id = ?", userID).
		Update("refresh_token", refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryDB) FindSession(userID uuid.UUID) (*domain.UserSession, error) {
	var userSession domain.UserSession
	if err := r.db.Where("user_id = ?", userID).First(&userSession).Error; err != nil {
		return nil, err
	}
	return &userSession, nil
}

func (r *UserRepositoryDB) CreateUserAndSession(user *domain.User, refreshToken *string) (*domain.User, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	userSession := domain.UserSession{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}
	if err := tx.Create(&userSession).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	user.UserSession = userSession
	return user, nil
}

func (r *UserRepositoryDB) FindUserWithRelation(ID uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload(clause.Associations).Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) FindAllUsersWithRelation() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Preload(clause.Associations).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryDB) FindUserBookmarks(userID uuid.UUID) ([]domain.Bookmark, error) {
	var bookmarks []domain.Bookmark
	if err := r.db.Preload("Post").Where("user_id = ?", userID).Find(&bookmarks).Error; err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (r *UserRepositoryDB) AddBookmark(bookmark domain.Bookmark) error {
	if err := r.db.Create(&bookmark).Error; err != nil {
		return err
	}
	return nil
}
