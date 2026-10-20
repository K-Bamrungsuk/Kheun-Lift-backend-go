package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"kheun-Lift-backend-go/internal/middleware"
)

func GetMe(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"message": "Please login.",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": session.User,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("Write user response failed: %v", err)
	}
}
