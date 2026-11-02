package weightclasses

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
		gender := c.Query("gender")
		if gender != "male" && gender != "female" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "gender must be male or female",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		weightClasses, err := list(ctx, pool, gender)
		if err != nil {
			log.Printf("List weight classes failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to load weight classes",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Get weight classes successfully",
			"data":    weightClasses,
		})
	}
}
