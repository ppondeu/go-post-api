package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/dto"
	errs "github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/logger"
	"github.com/ppondeu/go-post-api/internal/response"
	"github.com/ppondeu/go-post-api/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserService
	validator   *validator.Validate
}

func NewUserHandler(service usecase.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: service,
		validator:   validator,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	usersResponse := make([]dto.UserResponseDto, 0)
	for _, user := range users {
		usersResponse = append(usersResponse, dto.UserResponseDto{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			ShortBio: user.ShortBio,
		})
	}
	response.NewSuccessResponse(c, usersResponse)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequestError("id is invalid"))
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	userResponse := dto.UserResponseDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ShortBio: user.ShortBio,
	}
	response.NewSuccessResponse(c, userResponse)
}

func (h *UserHandler) GetUserSession(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequestError("id is invalid"))
		return
	}

	session, err := h.userService.GetUserSession(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	response.NewSuccessResponse(c, session)
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	userResponse := dto.UserResponseDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ShortBio: user.ShortBio,
	}
	response.NewSuccessResponse(c, userResponse)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	userResponse := dto.UserResponseDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ShortBio: user.ShortBio,
	}

	response.NewSuccessResponse(c, userResponse)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var craeteUserDto dto.CreateUserDto
	if err := c.ShouldBindJSON(&craeteUserDto); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(craeteUserDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, err)
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errs.NewBadRequestError(err.Error()))
		return
	}

	user, err := h.userService.CreateUserAndSession(&craeteUserDto, nil)
	if err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	userResponse := dto.UserResponseDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ShortBio: user.ShortBio,
	}
	response.NewCreatedResponse(c, userResponse)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	logger.Info(fmt.Sprintf("idParam: %v", idParam))
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequestError("id is invalid"))
		return
	}

	var updateUserDto dto.UpdateUserDto
	if err := c.ShouldBindJSON(&updateUserDto); err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(updateUserDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error(err)
			response.NewErrorResponse(c, err)
			return
		}
		logger.Error(err)
		response.NewErrorResponse(c, errs.NewBadRequestError(err.Error()))
		return
	}

	user, err := h.userService.UpdateUser(id, &updateUserDto)
	if err != nil {
		logger.Error(err)
		response.NewErrorResponse(c, err)
		return
	}

	userResponse := dto.UserResponseDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ShortBio: user.ShortBio,
	}
	response.NewSuccessResponse(c, userResponse)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequestError("id is invalid"))
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	response.NewSuccessResponse(c, nil)
}

// func (h *UserHandler) UpdateUserSession(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := uuid.Parse(idParam)
// 	if err != nil {
// 		response.NewErrorResponse(c, errs.NewBadRequestError("id is invalid"))
// 		return
// 	}
// 	type RefreshToken struct {
// 		RefreshToken *string `json:"refresh_token" validate:"omitempty"`
// 	}
// 	var body RefreshToken
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		logger.Error(err)
// 		response.NewErrorResponse(c, err)
// 		return
// 	}

// 	if err := h.validator.Struct(body); err != nil {
// 		if _, ok := err.(*validator.InvalidValidationError); ok {
// 			logger.Error(err)
// 			response.NewErrorResponse(c, err)
// 			return
// 		}
// 		logger.Error(err)
// 		response.NewErrorResponse(c, errs.NewBadRequestError(err.Error()))
// 		return
// 	}

// 	err = h.userService.UpdateUserSession(id, body.RefreshToken)
// 	if err != nil {
// 		response.NewErrorResponse(c, err)
// 		return
// 	}

// 	response.NewSuccessResponse(c, nil)
// }
