package usecase

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/domain"
	"github.com/ppondeu/go-post-api/internal/dto"
	"github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/logger"
	"gorm.io/gorm"
)

type PostService interface {
	GetAllPosts() ([]domain.Post, error)
	GetPostByID(ID uuid.UUID) (*domain.Post, error)
	GetPostsByUserID(userID uuid.UUID) ([]domain.Post, error)
	CreatePost(post dto.CreatePostDto) (*domain.Post, error)
	UpdatePost(ID uuid.UUID, post dto.UpdatePostDto) (*domain.Post, error)
	DeletePost(ID uuid.UUID) error

	GetAllTags() ([]domain.Tag, error)
}

type postServiceImpl struct {
	postRepo domain.PostRepository
}

func NewPostService(postRepo domain.PostRepository) PostService {
	return &postServiceImpl{
		postRepo: postRepo,
	}
}

func (p *postServiceImpl) GetAllPosts() ([]domain.Post, error) {
	posts, err := p.postRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *postServiceImpl) GetPostByID(ID uuid.UUID) (*domain.Post, error) {
	post, err := p.postRepo.FindByID(ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Post not found")
		}
		return nil, errors.NewBadRequestError(err.Error())
	}

	return post, nil
}

func (p *postServiceImpl) GetPostsByUserID(userID uuid.UUID) ([]domain.Post, error) {
	posts, err := p.postRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Posts not found")
		}
		return nil, errors.NewBadRequestError(err.Error())
	}

	return posts, nil
}

func (p *postServiceImpl) CreatePost(postDto dto.CreatePostDto) (*domain.Post, error) {

	newPost := domain.Post{
		Title:   postDto.Title,
		Content: postDto.Content,
		UserID:  postDto.UserID,
		Tags:    postDto.Tags,
	}

	post, err := p.postRepo.Save(newPost)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return post, nil
}

func (p *postServiceImpl) UpdatePost(ID uuid.UUID, postDto dto.UpdatePostDto) (*domain.Post, error) {
	updatePost := domain.Post{
		Title:   postDto.Title,
		Content: postDto.Content,
		Tags:    postDto.Tags,
	}

	post, err := p.postRepo.Update(ID, updatePost)
	if err != nil {
		logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Post not found")
		}
		return nil, err
	}

	return post, nil
}

func (p *postServiceImpl) DeletePost(ID uuid.UUID) error {
	err := p.postRepo.Delete(ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *postServiceImpl) GetAllTags() ([]domain.Tag, error) {
	tags, err := p.postRepo.FindAllTags()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tags, nil
}
