package handler

import (
	"log"
	"net/http"
	"time"

	"kheun-Lift-backend-go/internal/auth/service"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLogoutHandler(
	pool *pgxpool.Pool,
	secureCookie bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := middleware.SessionFromContext(r.Context())
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"message": "Please login.",
			})
			return
		}

		if err := service.Logout(
			r.Context(), pool, session.TokenHash,
		); err != nil {
			log.Print("Delete session failed")
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"message": "Unable to logout.",
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "kheun_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   secureCookie,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0).UTC(),
		})

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Logged out successfully.",
		})
	}
}
