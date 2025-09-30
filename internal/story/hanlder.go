package story

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LoaderService struct {
	repo *Repository
}

func NewLoaderService(repo *Repository) *LoaderService {
	return &LoaderService{
		repo: repo,
	}
}

func (s *LoaderService) LoadStoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var storiesData []StoryData
	if err := json.NewDecoder(r.Body).Decode(&storiesData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.repo.LoadStoriesFromData(r.Context(), storiesData); err != nil {
		log.Printf("Error loading stories: %v", err)
		http.Error(w, "Failed to load stories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Stories loaded successfully")
}
