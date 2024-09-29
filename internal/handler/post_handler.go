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
