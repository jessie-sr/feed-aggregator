package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jessie-sr/rss-aggregator/internal/db"
)

func (apiCig *apiConfig) handleCreateFeedSaved(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameter struct {
		FeedID string `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body) // Create a JSON decoder that reads from the request body

	params := parameter{}
	err := decoder.Decode(&params) // Read from the request body into the params struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprintln("Error parsing JSON:", err))
		return
	}

	// Create a new feed using DB.CreateFeed
	feed, err := apiCig.DB.CreateFeedSaved(r.Context(), db.CreateFeedSavedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %s", err))
		return
	}

	// Return the feed through our custom feed model
	respondWithJSON(w, 201, dbFeedToFeed(feed))
}

func (apiCig *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	// Get all the feeds using DB.GetFeeds
	feeds, err := apiCig.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feeds: %s", err))
		return
	}

	// Return the feeds through our custom feeds model
	respondWithJSON(w, 200, dbFeedsToFeeds(feeds))
}
