package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotFoundError struct {
	Message string `json:"message"`
}

func (n *NotFoundError) Error() string {
	return n.Message
}

type BadRequestError struct {
	Message string `json:"message"`
}

func (n *BadRequestError) Error() string {
	return n.Message
}

type UnauthorizedError struct {
	Message string `json:"message"`
}

func (n *UnauthorizedError) Error() string {
	return n.Message
}

func BuildErrorAPI(c *gin.Context, err error) {
	switch err.(type) {
	case *NotFoundError:
		c.AbortWithStatusJSON(
			http.StatusNotFound, map[string]interface{}{
				"message": "StatusNotFound",
				"status":  "failed",
				"error":   err.Error(),
			},
		)
		return
	case *BadRequestError:
		c.AbortWithStatusJSON(
			http.StatusBadRequest, map[string]interface{}{
				"message": "StatusBadRequest",
				"status":  "failed",
				"error":   err.Error(),
			},
		)
		return
	case *UnauthorizedError:
		c.AbortWithStatusJSON(
			http.StatusUnauthorized, map[string]interface{}{
				"message": "StatusUnauthorized",
				"status":  "failed",
				"error":   err.Error(),
			},
		)
		return
	case error:
		c.AbortWithStatusJSON(
			http.StatusInternalServerError, map[string]interface{}{
				"message": "StatusInternalServerError",
				"status":  "failed",
				"error":   "internal server error",
			},
		)
		return
	}
}
