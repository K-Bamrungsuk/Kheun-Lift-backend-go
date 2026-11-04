package routes

import (
	"net/http"
	"os"

	authhandler "kheun-Lift-backend-go/internal/auth/handler"
	"kheun-Lift-backend-go/internal/exercises"
	"kheun-Lift-backend-go/internal/lifts"
	appmiddleware "kheun-Lift-backend-go/internal/middleware"
	userhandler "kheun-Lift-backend-go/internal/user/handler"
	"kheun-Lift-backend-go/internal/weightclasses"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(router *gin.Engine, pool *pgxpool.Pool) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authRateLimit := appmiddleware.RateLimit()
	router.POST("/auth/register", authRateLimit, authhandler.NewRegisterHandler(pool))
	router.POST("/auth/login", authRateLimit, authhandler.NewLoginHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))
	router.GET("/exercises", exercises.NewExerciseHandler(pool))
	router.GET("/exercises/:id", exercises.NewGetExerciseByIDHandler(pool))

	crossOrigin := http.NewCrossOriginProtection()
	protected := router.Group("")
	protected.Use(func(c *gin.Context) {
		if err := crossOrigin.Check(c.Request); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Cross-origin request blocked.",
			})
			return
		}
		c.Next()
	}, appmiddleware.RequireSession(pool))

	protected.GET("/users/me", userhandler.NewGetMeHandler(pool))
	protected.GET("/weight-classes", weightclasses.NewWeightClassHandler(pool))
	protected.PATCH("/users/me", userhandler.NewPatchMeHandler(pool))
	protected.POST("/lifts", lifts.NewCreateHandler(pool))
	protected.GET("/lifts/me", lifts.NewListMineHandler(pool))
	protected.GET("/lifts/users/:id", lifts.NewListByUserHandler(pool))
	protected.GET("/lifts/:id", lifts.NewGetLiftByIDHandler(pool))
	protected.PATCH("/lifts/:id", lifts.NewUpdateHandler(pool))
	protected.DELETE("/lifts/:id", lifts.NewDeleteHandler(pool))
	protected.POST("/auth/logout", authhandler.NewLogoutHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))
}
