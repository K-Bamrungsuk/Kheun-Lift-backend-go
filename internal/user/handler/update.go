package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"kheun-Lift-backend-go/internal/middleware"
	"kheun-Lift-backend-go/internal/user/model"
	"kheun-Lift-backend-go/internal/user/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func decodeUpdateUser(body io.Reader) (model.UpdateUserInput, error) {
	var input model.UpdateUserInput
	var fields map[string]json.RawMessage

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&fields); err != nil {
		return input, errors.New("Invalid JSON or request body too large.")
	}
	if len(fields) == 0 {
		return input, errors.New("Provide at least one field to update.")
	}

	var dateOfBirth string
	for field, raw := range fields {
		var target any
		switch field {
		case "username":
			target = &input.Username
		case "profileImage":
			target = &input.ProfileImage
		case "bodyWeight":
			target = &input.BodyWeight
		case "height":
			target = &input.Height
		case "gender":
			target = &input.Gender
		case "dateOfBirth":
			target = &dateOfBirth
		default:
			return input, fmt.Errorf("Unsupported field: %s.", field)
		}

		if strings.TrimSpace(string(raw)) == "null" {
			return input, fmt.Errorf("%s cannot be null.", field)
		}
		if err := json.Unmarshal(raw, target); err != nil {
			return input, fmt.Errorf("Invalid value for %s.", field)
		}
	}

	if _, supplied := fields["dateOfBirth"]; supplied {
		date, err := time.Parse(time.DateOnly, dateOfBirth)
		if err != nil {
			date, err = time.Parse(time.RFC3339, dateOfBirth)
		}
		if err != nil {
			return input, errors.New("Date of birth must be YYYY-MM-DD or RFC3339.")
		}

		// Birthdays retain the supplied calendar date regardless of timezone.
		normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		input.DateOfBirth = &normalized
	}

	return input, nil
}

func validateUpdateUser(input *model.UpdateUserInput) error {
	if input.Username != nil {
		*input.Username = strings.TrimSpace(*input.Username)
		if n := utf8.RuneCountInString(*input.Username); n < 1 || n > 20 {
			return errors.New("Username must be between 1 and 20 characters.")
		}
	}
	if input.BodyWeight != nil && (*input.BodyWeight < 1 || *input.BodyWeight > 999) {
		return errors.New("Body weight must be between 1 and 999.")
	}
	if input.Height != nil && (*input.Height < 1 || *input.Height > 999) {
		return errors.New("Height must be between 1 and 999.")
	}
	if input.Gender != nil && *input.Gender != "male" && *input.Gender != "female" {
		return errors.New("Gender must be male or female.")
	}

	if input.DateOfBirth != nil {
		day := input.DateOfBirth.Format(time.DateOnly)
		if day < "1900-01-01" || day > time.Now().Format(time.DateOnly) {
			return errors.New("Date of birth must be between 1900-01-01 and today.")
		}
	}
	return nil
}

func NewPatchMeHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := middleware.SessionFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Please login.",
			})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
		input, err := decodeUpdateUser(c.Request.Body)
		if err == nil {
			err = validateUpdateUser(&input)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		user, err := service.UpdateUser(c.Request.Context(), pool, session.User.ID, input)
		if err != nil {
			var pgErr *pgconn.PgError
			switch {
			case errors.As(err, &pgErr) &&
				pgErr.Code == "23505" &&
				pgErr.ConstraintName == "User_username_key":
				c.JSON(http.StatusConflict, map[string]string{
					"message": "Username is already registered.",
				})
			case errors.Is(err, pgx.ErrNoRows):
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Session is invalid. Please login.",
				})
			default:
				log.Print("Update user failed")
				c.JSON(http.StatusInternalServerError, map[string]string{
					"message": "Unable to update user.",
				})
			}
			return
		}

		c.JSON(http.StatusOK, map[string]any{
			"user": user,
		})
	}
}
