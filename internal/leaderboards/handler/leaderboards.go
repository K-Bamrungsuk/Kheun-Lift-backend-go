package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"kheun-Lift-backend-go/internal/leaderboards/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func parseID(c *gin.Context, name string) (int32, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 32)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": name + " must be a positive 32-bit integer."})
		return 0, false
	}
	return int32(id), true
}

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
		exerciseID, ok := parseID(c, "exerciseId")
		if !ok {
			return
		}
		var weightClassID *int32
		if c.Param("weightClassId") != "" {
			id, ok := parseID(c, "weightClassId")
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
		c.JSON(http.StatusOK, records)
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
