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
	AddBookmark(userID, PostID uuid.UUID) error
	RemoveBookmark(userID, PostID uuid.UUID) error
	LikePost(userID, PostID uuid.UUID) error
	UnlikePost(userID, PostID uuid.UUID) error

	AddComment(createCommentDto dto.CreateCommentDto) (*domain.Comment, error)
	UpdateComment(commentID uuid.UUID, content string) (*domain.Comment, error)
	DeleteComment(commentID uuid.UUID) error
	GetCommentsByPost(postID uuid.UUID) ([]domain.Comment, error)
	GetCommentByID(commentID uuid.UUID) (*domain.Comment, error)
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

	if len(postDto.Tags) > 0 {
		postDto.Tags = []string{}
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

	if updatePost.Tags == nil {
		updatePost.Tags = []string{}
	} else {
		updatePost.Tags = postDto.Tags
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

func (p *postServiceImpl) AddBookmark(userID, PostID uuid.UUID) error {
	_, err := p.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(PostID)
	if err != nil {
		return err
	}

	if post.UserID == userID.String() {
		logger.Error("You can't bookmark your own post")
		return errors.NewBadRequestError("You can't bookmark your own post")
	}

	bookmark := domain.Bookmark{
		UserID: userID.String(),
		PostID: PostID.String(),
	}

	err = p.postRepo.AddBookmark(bookmark)
	if err != nil {
		return errors.NewBadRequestError("Bookmark already exists")
	}

	return nil
}

func (p *postServiceImpl) RemoveBookmark(userID, PostID uuid.UUID) error {
	_, err := p.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(PostID)
	if err != nil {
		return err
	}

	if post.UserID == userID.String() {
		logger.Error("You can't remove bookmark from your own post")
		return errors.NewBadRequestError("You can't remove bookmark from your own post")
	}

	err = p.postRepo.RemoveBookmark(userID, PostID)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *postServiceImpl) LikePost(userID, PostID uuid.UUID) error {
	_, err := p.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(PostID)
	if err != nil {
		return err
	}

	if post.UserID == userID.String() {
		logger.Error("You can't like your own post")
		return errors.NewBadRequestError("You can't like your own post")
	}

	like := domain.Like{
		UserID: userID.String(),
		PostID: PostID.String(),
	}

	err = p.postRepo.LikePost(like)
	if err != nil {
		logger.Error(err)
		return errors.NewBadRequestError("You already liked this post")
	}

	return nil
}

func (p *postServiceImpl) UnlikePost(userID, PostID uuid.UUID) error {
	_, err := p.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	post, err := p.GetPostByID(PostID)
	if err != nil {
		return err
	}

	if post.UserID == userID.String() {
		logger.Error("You can't unlike your own post")
		return errors.NewBadRequestError("You can't unlike your own post")
	}

	err = p.postRepo.UnlikePost(userID, PostID)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *postServiceImpl) AddComment(createCommentDto dto.CreateCommentDto) (*domain.Comment, error) {
	userID, err := uuid.Parse(createCommentDto.UserID)
	if err != nil {
		return nil, err
	}
	PostID, err := uuid.Parse(createCommentDto.PostID)
	if err != nil {
		return nil, err
	}

	_, err = p.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	_, err = p.GetPostByID(PostID)
	if err != nil {
		return nil, err
	}

	comment := domain.Comment{
		Content:  createCommentDto.Content,
		UserID:   createCommentDto.UserID,
		PostID:   createCommentDto.PostID,
		ParentID: createCommentDto.ParentID,
	}

	newComment, err := p.postRepo.AddComment(comment)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return newComment, nil
}

func (p *postServiceImpl) UpdateComment(commentID uuid.UUID, content string) (*domain.Comment, error) {

	comment := domain.Comment{
		Content: content,
	}

	updatedComment, err := p.postRepo.UpdateComment(commentID, comment)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updatedComment, nil
}

func (p *postServiceImpl) DeleteComment(commentID uuid.UUID) error {
	err := p.postRepo.DeleteComment(commentID)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *postServiceImpl) GetCommentsByPost(postID uuid.UUID) ([]domain.Comment, error) {
	comments, err := p.postRepo.FindCommentsByPostID(postID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var replies []domain.Comment
	for i := range comments {
		replies, err = p.postRepo.FindRepliesByCommentID(uuid.MustParse(comments[i].ID))
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		comments[i].Replies = replies
	}

	return comments, nil
}

func (p *postServiceImpl) GetCommentByID(commentID uuid.UUID) (*domain.Comment, error) {
	comment, err := p.postRepo.FindCommentByID(commentID)
	if err != nil {
		logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Comment not found")
		}
		return nil, errors.NewBadRequestError(err.Error())
	}

	var replies []domain.Comment
	replies, err = p.postRepo.FindRepliesByCommentID(commentID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	comment.Replies = replies

	return comment, nil
}
