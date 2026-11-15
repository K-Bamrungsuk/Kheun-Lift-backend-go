package handler

import (
	"errors"
	"net/http"
	"strings"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/auth/service"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLoginHandler(pool *pgxpool.Pool, cookie middleware.SessionCookie) gin.HandlerFunc {
	hashSlots := make(chan struct{}, 2)

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		fail := func(status int, message string) {
			c.JSON(status, map[string]string{
				"message": message,
			})
		}

		var input model.LoginInput
		if !middleware.DecodeBody(c, &input) {
			return
		}

		input.Email = strings.TrimSpace(input.Email)
		if input.Email == "" || len(input.Email) > 254 ||
			input.Password == "" || len(input.Password) > 1024 {
			fail(http.StatusBadRequest, "Email and password are required within allowed lengths.")
			return
		}

		select {
		case hashSlots <- struct{}{}:
			defer func() { <-hashSlots }()
		default:
			fail(http.StatusServiceUnavailable, "Login is busy. Please try again.")
			return
		}

		result, err := service.Login(c.Request.Context(), pool, input)
		if errors.Is(err, service.ErrInvalidCredentials) {
			fail(http.StatusUnauthorized, "Invalid email or password.")
			return
		}

		if err != nil {
			fail(http.StatusInternalServerError, "Unable to login.")
			return
		}

		cookie.Set(c.Writer, result.Token, result.ExpiresAt)
		c.JSON(http.StatusOK, map[string]any{
			"message": "Welcome! Ready to update your PRs?",
			"user":    result.User,
		})
	}
}
