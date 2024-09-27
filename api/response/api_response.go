package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/errs"
)

type ApiResponse struct {
	StatusCode uint16      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func NewApiResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, ApiResponse{
		StatusCode: uint16(statusCode),
		Message:    message,
		Data:       data,
	})
}

func NewSuccessResponse(c *gin.Context, data interface{}) {
	NewApiResponse(c, 200, "success", data)
}

func NewErrorResponse(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errs.AppError:
		NewApiResponse(c, e.Code, e.Message, nil)
	default:
		NewApiResponse(c, http.StatusInternalServerError, "An unexpected error occurred", nil)
	}
}

func NewCreatedResponse(c *gin.Context, data interface{}) {
	NewApiResponse(c, 201, "created", data)
}
