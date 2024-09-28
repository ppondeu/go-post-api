package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/response"
	"github.com/ppondeu/go-post-api/internal/usecase"
)

type Payload struct {
	Claims *usecase.UserClaims
	Token  string
}

func ValidateAccessToken(jwtService usecase.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("accessToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token provided"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			c.Abort()
			return
		}

		fmt.Println(claims)

		payload := Payload{
			Claims: claims,
			Token:  tokenString,
		}
		c.Set("payload", payload)

		c.Next()
	}
}

func ValidateRefreshToken(jwtService usecase.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("refreshToken")
		if err != nil {
			response.NewErrorResponse(c, errors.NewUnauthorizedError("No refresh token provided"))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString, "refresh")
		if err != nil {
			response.NewErrorResponse(c, errors.NewUnauthorizedError("Invalid refresh token"))
			c.Abort()
			return
		}

		payload := Payload{
			Claims: claims,
			Token:  tokenString,
		}
		c.Set("payload", payload)

		c.Next()
	}
}
