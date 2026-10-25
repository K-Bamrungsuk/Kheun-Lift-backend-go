package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/auth/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func validateRegister(input model.RegisterInput) string {
	if n := utf8.RuneCountInString(input.Username); n < 1 || n > 20 {
		return "Username must be between 1 and 20 characters."
	}

	address, err := mail.ParseAddress(input.Email)
	if err != nil || address.Address != input.Email ||
		len(input.Email) > 254 ||
		strings.ContainsAny(input.Email, "\r\n") {
		return "Please enter a valid email address."
	}

	if utf8.RuneCountInString(input.Password) < 6 {
		return "Password must be at least 6 characters."
	}
	if len(input.Password) > 1024 {
		return "Password must be at most 1024 bytes."
	}

	return ""
}

func NewRegisterHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	hashSlots := make(chan struct{}, 2)

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
		decoder := json.NewDecoder(c.Request.Body)

		var input model.RegisterInput
		if err := decoder.Decode(&input); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid JSON or request body too large.",
			})
			return
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Request body must contain one JSON object.",
			})
			return
		}

		input.Username = strings.TrimSpace(input.Username)
		input.Email = strings.TrimSpace(input.Email)

		if message := validateRegister(input); message != "" {
			c.JSON(http.StatusBadRequest, map[string]string{
				"message": message,
			})
			return
		}

		select {
		case hashSlots <- struct{}{}:
			defer func() { <-hashSlots }()
		default:
			c.JSON(http.StatusServiceUnavailable, map[string]string{
				"message": "Registration is busy. Please try again.",
			})
			return
		}

		user, err := service.RegisterUser(c.Request.Context(), pool, input)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, map[string]string{
					"message": "Username or email is already registered.",
				})
				return
			}

			log.Print("Create user failed")
			c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Unable to register.",
			})
			return
		}

		c.JSON(http.StatusCreated, map[string]any{
			"message": "Registration successful! Ready to smash your PR?",
			"user":    user,
		})
	}
}
