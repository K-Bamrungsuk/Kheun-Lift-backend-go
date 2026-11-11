package routes

import (
	"log"
	"net/http"
	"os"

	authhandler "kheun-Lift-backend-go/internal/auth/handler"
	exerciseshandler "kheun-Lift-backend-go/internal/exercises/handler"
	leaderboardshandler "kheun-Lift-backend-go/internal/leaderboards/handler"
	liftshandler "kheun-Lift-backend-go/internal/lifts/handler"
	appmiddleware "kheun-Lift-backend-go/internal/middleware"
	userhandler "kheun-Lift-backend-go/internal/user/handler"
	weightclasseshandler "kheun-Lift-backend-go/internal/weightclasses/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(router *gin.Engine, pool *pgxpool.Pool) {

	router.Use(appmiddleware.CORS(os.Getenv("APP_ORIGIN")))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public authentication
	authRateLimit := appmiddleware.RateLimit()
	router.POST("/auth/register", authRateLimit, authhandler.NewRegisterHandler(pool))
	router.POST("/auth/login", authRateLimit, authhandler.NewLoginHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))

	// Public exercises and leaderboards
	router.GET("/exercises", exerciseshandler.NewExerciseHandler(pool))
	router.GET("/exercises/:id", exerciseshandler.NewGetExerciseByIDHandler(pool))
	router.GET("/leaderboards/exercises/:exerciseId", leaderboardshandler.NewListHandler(pool))
	router.GET("/leaderboards/exercises/:exerciseId/weight-classes/:weightClassId", leaderboardshandler.NewListHandler(pool))
	router.GET("/leaderboards/random", leaderboardshandler.NewRandomHandler(pool))

	// Routes requiring a session
	crossOrigin := http.NewCrossOriginProtection()
	if origin := os.Getenv("APP_ORIGIN"); origin != "" {
		if err := crossOrigin.AddTrustedOrigin(origin); err != nil {
			log.Fatalf("invalid APP_ORIGIN: %v", err)
		}
	}

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

	// Profile and weight classes
	protected.GET("/users/me", userhandler.NewGetMeHandler(pool))
	protected.GET("/users/:id", userhandler.NewGetUserByIDHandler(pool))
	protected.GET("/weight-classes", weightclasseshandler.NewWeightClassHandler(pool))
	protected.PATCH("/users/me", userhandler.NewPatchMeHandler(pool))
	protected.DELETE("/users/me", userhandler.NewDeleteMeHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))

	// Lift records
	protected.POST("/lifts", liftshandler.NewCreateHandler(pool))
	protected.GET("/lifts/me", liftshandler.NewListMineHandler(pool))
	protected.GET("/lifts/users/:id", liftshandler.NewListByUserHandler(pool))
	protected.GET("/lifts/:id", liftshandler.NewGetLiftByIDHandler(pool))
	protected.PATCH("/lifts/:id", liftshandler.NewUpdateHandler(pool))
	protected.DELETE("/lifts/:id", liftshandler.NewDeleteHandler(pool))

	// Session
	protected.POST("/auth/logout", authhandler.NewLogoutHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))
}
