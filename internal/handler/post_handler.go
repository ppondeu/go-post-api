package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/dto"
	"github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/logger"
	"github.com/ppondeu/go-post-api/internal/response"
	"github.com/ppondeu/go-post-api/internal/usecase"
)

type PostHandler struct {
	postService usecase.PostService
	validator   *validator.Validate
}

func NewPostHandler(service usecase.PostService, validator *validator.Validate) *PostHandler {
	return &PostHandler{
		postService: service,
		validator:   validator,
	}
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	response.NewSuccessResponse(c, posts)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	post, err := h.postService.GetPostByID(postId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, post)
}

func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	posts, err := h.postService.GetPostsByUserID(userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, posts)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var createPostDto dto.CreatePostDto
	if err := c.ShouldBindJSON(&createPostDto); err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(createPostDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	post, err := h.postService.CreatePost(createPostDto)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	var updatePostDto dto.UpdatePostDto
	if err := c.ShouldBindJSON(&updatePostDto); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(updatePostDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, err)
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	post, err := h.postService.UpdatePost(postId, updatePostDto)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	err = h.postService.DeletePost(postId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	response.NewSuccessResponse(c, nil)
}

func (h *PostHandler) GetTags(c *gin.Context) {
	tags, err := h.postService.GetAllTags()
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	response.NewSuccessResponse(c, tags)
}

func (h *PostHandler) AddBookmark(c *gin.Context) {
	var createBookmarkDto dto.BookmarkDto
	if err := c.ShouldBindJSON(&createBookmarkDto); err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(createBookmarkDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, errors.NewBadRequestError("Invalid request"))
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	userId, err := uuid.Parse(createBookmarkDto.UserID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	postId, err := uuid.Parse(createBookmarkDto.PostID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	err = h.postService.AddBookmark(userId, postId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, nil)
}

func (h *PostHandler) RemoveBookmark(c *gin.Context) {
	var removeBookmarkDto dto.BookmarkDto
	if err := c.ShouldBindJSON(&removeBookmarkDto); err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(removeBookmarkDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, errors.NewBadRequestError("Invalid request"))
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	userId, err := uuid.Parse(removeBookmarkDto.UserID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	postId, err := uuid.Parse(removeBookmarkDto.PostID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	err = h.postService.RemoveBookmark(userId, postId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, nil)
}

func (h *PostHandler) LikePost(c *gin.Context) {
	var likeDto dto.LikeDto
	if err := c.ShouldBindJSON(&likeDto); err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(likeDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, errors.NewBadRequestError("Invalid request"))
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	userId, err := uuid.Parse(likeDto.UserID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	postId, err := uuid.Parse(likeDto.PostID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	err = h.postService.LikePost(userId, postId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, nil)
}

func (h *PostHandler) UnlikePost(c *gin.Context) {
	var likeDto dto.LikeDto
	if err := c.ShouldBindJSON(&likeDto); err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(likeDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, errors.NewBadRequestError("Invalid request"))
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	userId, err := uuid.Parse(likeDto.UserID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	postId, err := uuid.Parse(likeDto.PostID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	err = h.postService.UnlikePost(userId, postId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, nil)
}

func (h *PostHandler) AddComment(c *gin.Context) {
	var createCommentDto dto.CreateCommentDto
	if err := c.ShouldBindJSON(&createCommentDto); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(createCommentDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, errors.NewBadRequestError("Invalid request"))
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	comment, err := h.postService.AddComment(createCommentDto)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, comment)
}

func (h *PostHandler) UpdateComment(c *gin.Context) {
	ID := c.Param("id")
	commentID, err := uuid.Parse(ID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	var updateCommentDto dto.UpdateCommentDto
	if err := c.ShouldBindJSON(&updateCommentDto); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(updateCommentDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, err)
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	comment, err := h.postService.UpdateComment(commentID, updateCommentDto.Content)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, comment)
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
	ID := c.Param("id")
	commentID, err := uuid.Parse(ID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	err = h.postService.DeleteComment(commentID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	response.NewSuccessResponse(c, nil)
}

func (h *PostHandler) GetCommentsByPostID(c *gin.Context) {
	ID := c.Param("id")
	postID, err := uuid.Parse(ID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	comments, err := h.postService.GetCommentsByPost(postID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, comments)
}

func (h *PostHandler) GetCommentByID(c *gin.Context) {
	ID := c.Param("id")
	commentID, err := uuid.Parse(ID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	comment, err := h.postService.GetCommentByID(commentID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, comment)
}
