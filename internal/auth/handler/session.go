package handler

import (
	"log"
	"net/http"
	"time"

	"kheun-Lift-backend-go/internal/auth/service"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLogoutHandler(
	pool *pgxpool.Pool,
	secureCookie bool,
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

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "kheun_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   secureCookie,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0).UTC(),
		})

		c.JSON(http.StatusOK, map[string]string{
			"message": "Logged out successfully.",
		})
	}
}
