package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/dto"
	"github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/middleware"
	"github.com/ppondeu/go-post-api/internal/response"
	"github.com/ppondeu/go-post-api/internal/usecase"
)

type AuthHandler struct {
	authService usecase.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService usecase.AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{authService: authService, validator: validator}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var authRequestDto dto.AuthRequestDTO
	if err := c.ShouldBindJSON(&authRequestDto); err != nil {
		response.NewErrorResponse(c, errors.NewBadRequestError("invalid json"))
		return
	}

	if err := h.validator.Struct(authRequestDto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			response.NewErrorResponse(c, errors.NewBadRequestError("invalid json"))
			return
		}

		response.NewErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	tokenResponseDto, err := h.authService.Login(authRequestDto)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.SetCookie("accessToken", tokenResponseDto.AccessToken, 60*5, "/", "", false, true)
	c.SetCookie("refreshToken", tokenResponseDto.RefreshToken, 60*9, "/", "", false, true)
	response.NewSuccessResponse(c, tokenResponseDto)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	payload := c.MustGet("payload").(middleware.Payload)
	fmt.Printf("payload: %v\n", payload)

	userId, err := uuid.Parse(payload.Claims.Sub)
	if err != nil {
		response.NewErrorResponse(c, errors.NewBadRequestError("invalid user id"))
		return
	}

	err = h.authService.Logout(payload.Token, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.SetCookie("accessToken", "", -1, "/", "", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)
	response.NewSuccessResponse(c, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {

	payload := c.MustGet("payload").(middleware.Payload)
	fmt.Printf("payload: %v\n", payload)

	userId, err := uuid.Parse(payload.Claims.Sub)
	if err != nil {
		response.NewErrorResponse(c, errors.NewBadRequestError("invalid user id"))
		return
	}

	tokenResponseDto, err := h.authService.RefreshToken(payload.Token, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	fmt.Println(tokenResponseDto)
	c.SetCookie("accessToken", tokenResponseDto.AccessToken, 60*5, "/", "", false, true)
	c.SetCookie("refreshToken", tokenResponseDto.RefreshToken, 60*9, "/", "", false, true)
	response.NewSuccessResponse(c, tokenResponseDto)
}
