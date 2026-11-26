package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"kheun-Lift-backend-go/internal/middleware"
	"kheun-Lift-backend-go/internal/user/repository"
	"kheun-Lift-backend-go/internal/user/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDeleteMeHandler(pool *pgxpool.Pool, cookie middleware.SessionCookie, store *repository.ProfileImageStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := middleware.SessionFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login."})
			return
		}

		if err := service.DeleteUser(c.Request.Context(), pool, session.User.ID); err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Session is invalid. Please login."})
			default:
				log.Printf("Delete user failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete account."})
			}
			return
		}

		// The account is already gone, so a failed image delete is only logged.
		if session.User.ProfileImage != nil {
			ctx, cancel := context.WithTimeout(context.WithoutCancel(c.Request.Context()), 15*time.Second)
			defer cancel()
			if err := store.Delete(ctx, session.User.ID); err != nil {
				log.Printf("Delete profile image failed: %v", err)
			}
		}

		cookie.Clear(c.Writer)
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
	}
}
