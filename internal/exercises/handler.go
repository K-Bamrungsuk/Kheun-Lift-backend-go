package exercises

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func NewGetByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "id must be a positive 32-bit integer",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		exercise, err := getByID(ctx, pool, int32(id))
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Exercise not found",
			})
			return
		}
		if err != nil {
			log.Printf("Get exercise failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to load exercise",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"exercise": exercise,
		})
	}
}
