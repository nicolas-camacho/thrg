package token

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicolas-camacho/thrg/internal/contextutil"
	"github.com/nicolas-camacho/thrg/internal/core"
)

type UserLookup interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*core.UserLookupModel, error)
}

func GenerateTokenHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		adminID, ok := contextutil.GetUserIDFromContext(ctx)

		if !ok || adminID == uuid.Nil {
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

func ListTokensHandler(tokenRepo *Repository, userLookup UserLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokens, err := tokenRepo.GetAllTokens(r.Context())
		if err != nil {
			log.Printf("Error retrieving tokens: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tokenDTOs := make([]TokenDTO, 0, len(tokens))
		for _, t := range tokens {
			dto := TokenDTO{
				Value:     t.Value,
				IsUsed:    t.IsUsed,
				CreatedAt: t.CreatedAt,
			}

			if t.IsUsed && t.UsedByID != nil {
				playerModel, lookupErr := userLookup.GetUserByID(r.Context(), *t.UsedByID)
				if lookupErr != nil {
					log.Printf("Error retrieving user for token %s: %v", t.Value, lookupErr)
					dto.UsedByUsername = "Error fetching user"
				} else if playerModel != nil {
					dto.UsedByUsername = playerModel.Username
				} else {
					dto.UsedByUsername = "Unknown"
				}
			} else {
				dto.UsedByUsername = "N/A"
			}
			tokenDTOs = append(tokenDTOs, dto)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(tokenDTOs); err != nil {
			log.Printf("Error encoding tokens to JSON: %v", err)
		}
	}
}
