package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jessie-sr/rss-aggregator/internal/db"
)

func (apiCig *apiConfig) handleCreateFeedSaved(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body) // Create a JSON decoder that reads from the request body

	params := parameter{}
	err := decoder.Decode(&params) // Read from the request body into the params struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprintln("Error parsing JSON:", err))
		return
	}

	// Create a new user-feed relation using DB.CreateFeed
	saved, err := apiCig.DB.CreateFeedSaved(r.Context(), db.CreateFeedSavedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error saving feed: %s", err))
		return
	}

	// Return the user-feed relation through our custom feed_saved model
	respondWithJSON(w, 201, dbFeedSavedToFeedSaved(saved))
}

func (apiCig *apiConfig) handleGetSavedFeeds(w http.ResponseWriter, r *http.Request, user db.User) {
	// Get all the feeds saved by the user using DB.GetSavedFeeds
	saved_feeds, err := apiCig.DB.GetSavedFeeds(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting saved feeds: %s", err))
		return
	}

	// Return the saved feeds through our custom saved_feeds model
	respondWithJSON(w, 200, dbSavedFeedsToSavedFeeds(saved_feeds))
}

func (apiCig *apiConfig) handleUnsaveFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	// Extract the "feed_saved_id" parameter from the URL and parse it as a UUID
	id_str := chi.URLParam(r, "feed_saved_id")
	feed_saved_id, err := uuid.Parse(id_str)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintln("Error parsing JSON:", err))
		return
	}

	err = apiCig.DB.UnsaveFeed(r.Context(), db.UnsaveFeedParams{
		ID:     feed_saved_id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error unsaving feed: %s", err))
		return
	}

	respondWithJSON(w, 200, "Successfully unsaved the feed")
}
