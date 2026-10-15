package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"kheun-Lift-backend-go/internal/exercises"

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

	router.Get("/exercises", exercises.NewHandler(pool))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", port)
	return server.ListenAndServe()
}
