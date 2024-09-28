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

type FollowHandler struct {
	followService usecase.FollowService
	validator     *validator.Validate
}

func NewFollowHandler(followService usecase.FollowService, validator *validator.Validate) *FollowHandler {
	return &FollowHandler{
		followService: followService,
		validator:     validator,
	}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	var followeBody dto.FollowRequestDto
	if err := c.ShouldBindJSON(&followeBody); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	followerID, err := uuid.Parse(followeBody.FollowerID)
	if err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError("follower_id is invalid"))
		return
	}

	followedID, err := uuid.Parse(followeBody.FollowedID)
	if err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError("followed_id is invalid"))
		return
	}

	err = h.followService.Follow(followerID, followedID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, nil)
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	var followeBody dto.FollowRequestDto
	if err := c.ShouldBindJSON(&followeBody); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	followerID, err := uuid.Parse(followeBody.FollowerID)
	if err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError("follower_id is invalid"))
		return
	}

	followedID, err := uuid.Parse(followeBody.FollowedID)
	if err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, errors.NewBadRequestError("followed_id is invalid"))
		return
	}

	err = h.followService.Unfollow(followerID, followedID)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, nil)
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.NewErrorResponse(c, errors.NewBadRequestError("id is invalid"))
		return
	}

	followers, err := h.followService.GetFollowers(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, followers)
}

func (h *FollowHandler) GetFollowedUsers(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.NewErrorResponse(c, errors.NewBadRequestError("id is invalid"))
		return
	}

	followedUsers, err := h.followService.GetFollowedUsers(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	response.NewSuccessResponse(c, followedUsers)
}
