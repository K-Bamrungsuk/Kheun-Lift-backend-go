package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/auth/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RequireSession(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")

			fail := func(status int, message string) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)

				if err := json.NewEncoder(w).Encode(map[string]string{
					"message": message,
				}); err != nil {
					log.Printf("Write auth error response failed: %v", err)
				}
			}

			cookie, err := r.Cookie("kheun_session")
			if err != nil || cookie.Value == "" {
				fail(http.StatusUnauthorized, "Please login.")
				return
			}

			session, err := service.FindSession(r.Context(), pool, cookie.Value)

			if errors.Is(err, pgx.ErrNoRows) {
				fail(
					http.StatusUnauthorized,
					"Session is invalid or expired. Please login.",
				)
				return
			}
			if err != nil {
				log.Print("Find session failed")
				fail(http.StatusInternalServerError, "Unable to check session.")
				return
			}

			requestContext := WithSession(r.Context(), session)
			next.ServeHTTP(w, r.WithContext(requestContext))
		})
	}
}

type sessionContextKey struct{}

func WithSession(ctx context.Context, session model.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey{}, session)
}

func SessionFromContext(ctx context.Context) (model.Session, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(model.Session)
	return session, ok
}
