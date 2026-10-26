package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	authhandler "kheun-Lift-backend-go/internal/auth/handler"
	"kheun-Lift-backend-go/internal/exercises"
	appmiddleware "kheun-Lift-backend-go/internal/middleware"
	userhandler "kheun-Lift-backend-go/internal/user/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return errors.New("invalid database configuration; check DATABASE_URL")
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return errors.New("database connection failed; check credentials, SSL settings, and network")
	}
	cancel()
	log.Println("Database connected")

	port := os.Getenv("PORT")

	if port == "" {
		port = "3001"
	}

	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = false
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authRateLimit := appmiddleware.RateLimit()
	router.POST("/auth/register", authRateLimit, authhandler.NewRegisterHandler(pool))
	router.POST("/auth/login", authRateLimit, authhandler.NewLoginHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))
	router.GET("/exercises", exercises.NewHandler(pool))

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

	protected.GET("/users/me", userhandler.GetMe)
	protected.PATCH("/users/me", userhandler.NewPatchMeHandler(pool))
	protected.POST("/auth/logout", authhandler.NewLogoutHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", port)
	return server.ListenAndServe()
}
