package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONResponse(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	JSONResponse(c, statusCode, true, message, data)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	JSONResponse(c, statusCode, false, message, nil)
}
