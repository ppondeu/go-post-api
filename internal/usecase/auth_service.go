package usecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/dto"
	"github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/logger"
	"github.com/ppondeu/go-post-api/internal/utils"
)

type AuthService interface {
	Login(authRequestDto dto.AuthRequestDTO) (*dto.TokenResponseDto, error)
	RefreshToken(refreshToken string, ID uuid.UUID) (*dto.TokenResponseDto, error)
	Logout(refreshToken string, ID uuid.UUID) error
}

type authServiceImpl struct {
	userService UserService
	jwtService  JwtService
}

type UserClaims struct {
	jwt.RegisteredClaims
	Sub       string `json:"sub"`
	Username  string `json:"username"`
	TokenType string `json:"tokenType"`
}

func NewAuthService(userService UserService, jwtService JwtService) AuthService {
	return &authServiceImpl{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (s *authServiceImpl) Login(authRequestDto dto.AuthRequestDTO) (*dto.TokenResponseDto, error) {
	user, err := s.userService.GetUserByEmail(authRequestDto.Email)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = utils.CompareHashAndPassword(user.Password, authRequestDto.Password)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Sub:       user.ID,
		Username:  user.Username,
		TokenType: "access",
	}

	access, err := s.jwtService.GenerateToken(userClaims, "access")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userClaims.TokenType = "refresh"

	refresh, err := s.jwtService.GenerateToken(userClaims, "refresh")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = s.userService.UpdateUserSession(userID, refresh)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	authResponse := &dto.TokenResponseDto{
		AccessToken:  *access,
		RefreshToken: *refresh,
	}
	return authResponse, nil
}

func (s *authServiceImpl) RefreshToken(refreshToken string, ID uuid.UUID) (*dto.TokenResponseDto, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken, "refresh")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.NewForbiddenError("invalid token type")
	}

	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if user.UserSession.RefreshToken == nil || claims.Username != user.Username {
		logger.Error(err)
		return nil, errors.NewForbiddenError("invalid refresh token")
	}

	if *user.UserSession.RefreshToken != refreshToken {
		logger.Error(err)
		return nil, errors.NewForbiddenError("invalid refresh token")
	}

	userClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Sub:       user.ID,
		Username:  user.Username,
		TokenType: "access",
	}

	access, err := s.jwtService.GenerateToken(userClaims, "access")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userClaims.TokenType = "refresh"

	refresh, err := s.jwtService.GenerateToken(userClaims, "refresh")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = s.userService.UpdateUserSession(userID, refresh)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	authResponse := &dto.TokenResponseDto{
		AccessToken:  *access,
		RefreshToken: *refresh,
	}
	return authResponse, nil
}

func (s *authServiceImpl) Logout(refreshToken string, ID uuid.UUID) error {

	user, err := s.userService.GetUserSession(ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if user.RefreshToken == nil {
		logger.Error(err)
		return errors.NewForbiddenError("invalid refresh token")
	}

	if *user.RefreshToken != refreshToken {
		logger.Error(err)
		return errors.NewForbiddenError("invalid refresh token")
	}

	err = s.userService.UpdateUserSession(ID, nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil

}
