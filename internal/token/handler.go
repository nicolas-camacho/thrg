package token

import (
	"encoding/json"
	"log"
	"net/http"
)

func GenerateTokenHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		adminIDVal := r.Context().Value("user_id")
		adminID, ok := adminIDVal.(uint)

		if !ok || adminID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenValue, err := repo.CreateNewToken(ctx, adminID)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"token":   tokenValue,
			"message": "Token generated successfully shared with a new player.",
		})
	}
}
