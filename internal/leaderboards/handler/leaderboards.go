package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"kheun-Lift-backend-go/internal/leaderboards/repository"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrExerciseNotFound),
		errors.Is(err, repository.ErrWeightClassNotFound),
		errors.Is(err, repository.ErrWeightClassGender):
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	default:
		log.Printf("Leaderboard request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get leaderboard."})
	}
}

func NewListHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseID, ok := middleware.ParseID(c, "exerciseId")
		if !ok {
			return
		}
		var weightClassID *int32
		if c.Param("weightClassId") != "" {
			id, ok := middleware.ParseID(c, "weightClassId")
			if !ok {
				return
			}
			weightClassID = &id
		}

		query, err := url.ParseQuery(c.Request.URL.RawQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query string."})
			return
		}
		var gender *string
		if values, present := query["gender"]; present {
			if len(values) != 1 || (values[0] != "male" && values[0] != "female") {
				c.JSON(http.StatusBadRequest, gin.H{"message": "gender must be male or female."})
				return
			}
			gender = &values[0]
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		if err := repository.ValidateFilters(ctx, pool, exerciseID, weightClassID, gender); err != nil {
			respondError(c, err)
			return
		}
		records, err := repository.List(ctx, pool, exerciseID, weightClassID, gender)
		if err != nil {
			respondError(c, err)
			return
		}
		message := "No leaderboard available"
		if len(records) > 0 {
			message = "Get leaderboard successfully"
		}
		c.JSON(http.StatusOK, gin.H{"message": message, "data": records})
	}
}

func NewRandomHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		result, err := repository.Random(ctx, pool)
		if err != nil {
			respondError(c, err)
			return
		}

		message := "No leaderboard available"
		if result != nil {
			message = "Get random leaderboard successfully"
		}
		c.JSON(http.StatusOK, gin.H{"message": message, "data": result})
	}
}
