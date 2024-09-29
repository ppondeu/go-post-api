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
	AddBookmark(userId, postId uuid.UUID) error
	RemoveBookmark(userId, postId uuid.UUID) error
	LikePost(userId, postId uuid.UUID) error
	UnlikePost(userId, postId uuid.UUID) error
}

type postServiceImpl struct {
	postRepo    domain.PostRepository
	userService UserService
}

func NewPostService(postRepo domain.PostRepository, userService UserService) PostService {
	return &postServiceImpl{
		postRepo:    postRepo,
		userService: userService,
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
		return nil, err
	}

	return tags, nil
}

func (p *postServiceImpl) AddBookmark(userId, postId uuid.UUID) error {
	_, err := p.userService.GetUserByID(userId)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(postId)
	if err != nil {
		return err
	}

	if post.UserID == userId.String() {
		logger.Error("You can't bookmark your own post")
		return errors.NewBadRequestError("You can't bookmark your own post")
	}

	bookmark := domain.Bookmark{
		UserID: userId.String(),
		PostID: postId.String(),
	}

	err = p.postRepo.AddBookmark(bookmark)
	if err != nil {
		return errors.NewBadRequestError("Bookmark already exists")
	}

	return nil
}

func (p *postServiceImpl) RemoveBookmark(userId, postId uuid.UUID) error {
	_, err := p.userService.GetUserByID(userId)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(postId)
	if err != nil {
		return err
	}

	if post.UserID == userId.String() {
		logger.Error("You can't remove bookmark from your own post")
		return errors.NewBadRequestError("You can't remove bookmark from your own post")
	}

	err = p.postRepo.RemoveBookmark(userId, postId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *postServiceImpl) LikePost(userId, postId uuid.UUID) error {
	_, err := p.userService.GetUserByID(userId)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(postId)
	if err != nil {
		return err
	}

	if post.UserID == userId.String() {
		logger.Error("You can't like your own post")
		return errors.NewBadRequestError("You can't like your own post")
	}

	like := domain.Like{
		UserID: userId.String(),
		PostID: postId.String(),
	}

	err = p.postRepo.LikePost(like)
	if err != nil {
		logger.Error(err)
		return errors.NewBadRequestError("You already liked this post")
	}

	return nil
}

func (p *postServiceImpl) UnlikePost(userId, postId uuid.UUID) error {
	_, err := p.userService.GetUserByID(userId)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(postId)
	if err != nil {
		return err
	}

	if post.UserID == userId.String() {
		logger.Error("You can't unlike your own post")
		return errors.NewBadRequestError("You can't unlike your own post")
	}

	err = p.postRepo.UnlikePost(userId, postId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
