package handler

import (
	"log"
	"net/http"

	"kheun-Lift-backend-go/internal/auth/service"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLogoutHandler(
	pool *pgxpool.Pool,
	cookie middleware.SessionCookie,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := middleware.SessionFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Please login.",
			})
			return
		}

		if err := service.Logout(
			c.Request.Context(), pool, session.TokenHash,
		); err != nil {
			log.Print("Delete session failed")
			c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Unable to logout.",
			})
			return
		}

		cookie.Clear(c.Writer)

		c.JSON(http.StatusOK, map[string]string{
			"message": "Logged out successfully.",
		})
	}
}
