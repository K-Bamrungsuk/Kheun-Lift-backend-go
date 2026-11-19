package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"kheun-Lift-backend-go/internal/exercises/repository"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewExerciseHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		exercises, err := repository.List(ctx, pool)

		if err != nil {
			log.Printf("List exercises failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to load exercises",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Get exercises successfully",
			"data":    exercises,
		})
	}
}

func NewGetExerciseByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := middleware.ParseID(c, "id")
		if !ok {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		exercise, err := repository.GetByID(ctx, pool, id)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "model.Exercise not found",
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
			"message": "Get exercise successfully",
			"data":    exercise,
		})
	}
}
