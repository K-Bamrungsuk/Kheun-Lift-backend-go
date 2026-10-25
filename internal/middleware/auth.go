package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/auth/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RequireSession(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		fail := func(status int, message string) {
			c.AbortWithStatusJSON(status, gin.H{"message": message})
		}

		cookie, err := c.Request.Cookie("kheun_session")
		if err != nil || cookie.Value == "" {
			fail(http.StatusUnauthorized, "Please login.")
			return
		}

		session, err := service.FindSession(
			c.Request.Context(), pool, cookie.Value,
		)
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

		requestContext := WithSession(c.Request.Context(), session)
		c.Request = c.Request.WithContext(requestContext)
		c.Next()
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
