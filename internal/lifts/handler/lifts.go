package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"kheun-Lift-backend-go/internal/lifts/model"
	"kheun-Lift-backend-go/internal/lifts/repository"
	"kheun-Lift-backend-go/internal/lifts/service"
	"kheun-Lift-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func decodeBody(c *gin.Context, input any) bool {
	contentType, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if err != nil || contentType != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "Use application/json."})
		return false
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON, unsupported fields or request body too large."})
		return false
	}
	return true
}

func parseID(c *gin.Context) (int32, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id must be a positive 32-bit integer."})
		return 0, false
	}
	return int32(id), true
}

func sessionUserID(c *gin.Context) (int, bool) {
	session, ok := middleware.SessionFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login."})
		return 0, false
	}
	return session.User.ID, true
}

func respondError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "Unable to process lift record."
	switch {
	case errors.Is(err, service.ErrDuplicate):
		status, message = http.StatusConflict, service.ErrDuplicate.Error()
	case errors.Is(err, service.ErrProfile):
		status, message = http.StatusBadRequest, service.ErrProfile.Error()
	case errors.Is(err, service.ErrExercise):
		status, message = http.StatusNotFound, service.ErrExercise.Error()
	case errors.Is(err, service.ErrWeightClass):
		status, message = http.StatusNotFound, service.ErrWeightClass.Error()
	case errors.Is(err, pgx.ErrNoRows):
		status, message = http.StatusNotFound, "Lift record not found."
	default:
		log.Printf("Lift record request failed: %v", err)
	}
	c.JSON(status, gin.H{"message": message})
}

func NewCreateHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := sessionUserID(c)
		if !ok {
			return
		}
		var input model.CreateInput
		if !decodeBody(c, &input) {
			return
		}
		if input.ExerciseID <= 0 || input.Reps <= 0 || input.Weight <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "exerciseId and reps must be positive 32-bit integers; weight must be a positive finite number."})
			return
		}
		if input.VideoURL != nil {
			videoURL, err := url.Parse(*input.VideoURL)
			if err != nil || videoURL.Hostname() == "" || (videoURL.Scheme != "http" && videoURL.Scheme != "https") {
				c.JSON(http.StatusBadRequest, gin.H{"message": "videoUrl must be an absolute HTTP or HTTPS URL."})
				return
			}
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		record, err := service.Create(ctx, pool, userID, input)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Lift Record created successfully!", "data": record})
	}
}

func NewListMineHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := sessionUserID(c)
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		records, err := repository.ListByUser(ctx, pool, userID)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Get all lift records successfully", "data": records})
	}
}

func NewListByUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		records, err := repository.ListByUser(ctx, pool, int(id))
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Get user all lift records successfully", "data": records})
	}
}

func NewGetLiftByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		record, err := repository.GetByID(ctx, pool, id)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Get lift record successfully", "data": record})
	}
}

func NewUpdateHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := sessionUserID(c)
		if !ok {
			return
		}
		id, ok := parseID(c)
		if !ok {
			return
		}
		var input model.UpdateInput
		if !decodeBody(c, &input) {
			return
		}
		if input.Caption == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "caption is required and must be a string."})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		record, err := repository.UpdateCaption(ctx, pool, id, userID, *input.Caption)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Update lift record successfully", "data": record})
	}
}

func NewDeleteHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := sessionUserID(c)
		if !ok {
			return
		}
		id, ok := parseID(c)
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		if err := repository.DeleteByID(ctx, pool, id, userID); err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Delete lift record successfully"})
	}
}
