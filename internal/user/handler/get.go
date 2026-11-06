package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	authmodel "kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/middleware"
	"kheun-Lift-backend-go/internal/user/repository"
	weightclassmodel "kheun-Lift-backend-go/internal/weightclasses/model"
	weightclassrepository "kheun-Lift-backend-go/internal/weightclasses/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewGetUserByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		user, err := repository.GetByID(ctx, pool, int32(id))
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found.",
			})
			return
		}
		if err != nil {
			log.Printf("Get user failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to load profile.",
			})
			return
		}

		user.WeightClass, err = weightclassrepository.FindForUser(
			ctx, pool, user.Gender, user.BodyWeight,
		)
		if err != nil {
			log.Printf("Find user weight class failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to load profile.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func NewGetMeHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := middleware.SessionFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Please login.",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		weightClass, err := weightclassrepository.FindForUser(ctx, pool, session.User.Gender, session.User.BodyWeight)
		if err != nil {
			log.Printf("Find user weight class failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to load profile."})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user": struct {
				authmodel.User
				WeightClass *weightclassmodel.WeightClass `json:"weightClass"`
			}{User: session.User, WeightClass: weightClass},
		})
	}
}
