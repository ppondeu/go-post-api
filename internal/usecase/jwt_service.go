package usecase

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ppondeu/go-post-api/internal/logger"
)

type JwtService interface {
	GenerateToken(userClaims UserClaims, typeToken string) (*string, error)
	ValidateToken(tokenString string, typeToken string) (*UserClaims, error)
	GetAccessSecret() []byte
	GetRefreshSecret() []byte
}

type jwtServiceImpl struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJwtService(accessSecret, refreshSecret []byte) JwtService {
	return &jwtServiceImpl{accessSecret, refreshSecret}
}

func (s *jwtServiceImpl) GenerateToken(userClaims UserClaims, typeToken string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	var secret []byte
	if typeToken == "access" {
		secret = s.accessSecret
	} else if typeToken == "refresh" {
		secret = s.refreshSecret
	} else {
		return nil, fmt.Errorf("invalid token type")
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &tokenString, nil
}

func (s *jwtServiceImpl) ValidateToken(tokenString string, typeToken string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if typeToken == "access" {
			return s.accessSecret, nil
		} else if typeToken == "refresh" {
			return s.refreshSecret, nil
		} else {
			return nil, fmt.Errorf("invalid token type")
		}
	})

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.TokenType != typeToken {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

func (s *jwtServiceImpl) GetAccessSecret() []byte {
	return s.accessSecret
}

func (s *jwtServiceImpl) GetRefreshSecret() []byte {
	return s.refreshSecret
}
