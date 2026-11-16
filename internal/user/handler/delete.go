package handler

import (
	"errors"
	"log"
	"net/http"

	"kheun-Lift-backend-go/internal/middleware"
	"kheun-Lift-backend-go/internal/user/repository"
	"kheun-Lift-backend-go/internal/user/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
			var pgErr *pgconn.PgError
			switch {
			case errors.As(err, &pgErr) && pgErr.Code == "23503" &&
				pgErr.ConstraintName == "LiftRecord_userId_fkey":
				c.JSON(http.StatusConflict, gin.H{"message": "Delete your lift records before deleting your account."})
			case errors.Is(err, pgx.ErrNoRows):
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Session is invalid. Please login."})
			default:
				log.Print("Delete user failed")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete account."})
			}
			return
		}

		// The account is already gone, so a failed image delete is only logged.
		if session.User.ProfileImage != nil {
			if err := store.Delete(c.Request.Context(), session.User.ID); err != nil {
				log.Printf("Delete profile image failed: %v", err)
			}
		}

		cookie.Clear(c.Writer)
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
	}
}
