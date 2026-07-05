package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response structure.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
	})
}

func Fail(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:    -1,
		Message: msg,
	})
}

func FailWithData(c *gin.Context, httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    -1,
		Message: msg,
		Data:    data,
	})
}
