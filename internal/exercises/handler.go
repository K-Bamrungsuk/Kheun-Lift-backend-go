package exercises

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		exercises, err := list(ctx, pool)

		if err != nil {
			log.Printf("List exercises failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to load exercises",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"exercises": exercises,
		})
	}
}
