package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	authmodel "kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/middleware"
	weightclassmodel "kheun-Lift-backend-go/internal/weightclasses/model"
	weightclassrepository "kheun-Lift-backend-go/internal/weightclasses/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
