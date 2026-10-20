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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	router.Post("/auth/register", authhandler.NewRegisterHandler(pool))
	router.Post("/auth/login", authhandler.NewLoginHandler(
		pool,
		os.Getenv("COOKIE_SECURE") != "false",
	))
	router.Get("/exercises", exercises.NewHandler(pool))

	router.Group(func(protected chi.Router) {
		protected.Use(http.NewCrossOriginProtection().Handler)
		protected.Use(appmiddleware.RequireSession(pool))

		protected.Get("/users/me", userhandler.GetMe)
		protected.Patch("/users/me", userhandler.NewPatchMeHandler(pool))
		protected.Post("/auth/logout", authhandler.NewLogoutHandler(
			pool,
			os.Getenv("COOKIE_SECURE") != "false",
		))
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", port)
	return server.ListenAndServe()
}
