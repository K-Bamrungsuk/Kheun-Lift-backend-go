package handler

import (
	"net/http"

	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	session, ok := middleware.SessionFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Please login.",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"user": session.User,
	})
}
