package handler

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/auth/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLoginHandler(
	pool *pgxpool.Pool,
	secureCookie bool,
) http.HandlerFunc {
	hashSlots := make(chan struct{}, 2)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		fail := func(status int, message string) {
			writeJSON(w, status, map[string]string{
				"message": message,
			})
		}

		contentType, _, err := mime.ParseMediaType(
			r.Header.Get("Content-Type"),
		)
		if err != nil || contentType != "application/json" {
			fail(http.StatusUnsupportedMediaType, "Use application/json.")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
		decoder := json.NewDecoder(r.Body)

		var input model.LoginInput
		if err := decoder.Decode(&input); err != nil {
			fail(http.StatusBadRequest, "Invalid JSON or request body too large.")
			return
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			fail(http.StatusBadRequest, "Request body must contain one JSON object.")
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

		result, err := service.Login(r.Context(), pool, input)
		if errors.Is(err, service.ErrInvalidCredentials) {
			fail(http.StatusUnauthorized, "Invalid email or password.")
			return
		}

		if err != nil {
			fail(http.StatusInternalServerError, "Unable to login.")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "kheun_session",
			Value:    result.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   secureCookie,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
			Expires:  result.ExpiresAt,
		})
		writeJSON(w, http.StatusOK, map[string]any{
			"message": "Welcome! Ready to update your PRs?",
			"user":    result.User,
		})
	}
}
