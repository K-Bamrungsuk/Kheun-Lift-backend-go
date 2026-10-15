package exercises

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		exercises, err := list(ctx, pool)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			log.Printf("List exercises failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			if err := json.NewEncoder(w).Encode(map[string]string{
				"message": "Unable to load exercises",
			}); err != nil {
				log.Printf("Write exercises error response failed: %v", err)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]any{
			"exercises": exercises,
		}); err != nil {
			log.Printf("Write exercises response failed: %v", err)
		}
	}
}
