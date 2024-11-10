package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jessie-sr/rss-aggregator/internal/db"
)

func (apiCig *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body) // Creates a JSON decoder that reads from the request body

	params := parameter{}
	err := decoder.Decode(&params) // Reads from the request body into the params struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprintln("Error parsing JSON:", err))
		return
	}

	user, err := apiCig.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %s", err))
		return
	}

	respondWithJSON(w, 200, dbUserToUser(user))
}

func (apiCig *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {

}
