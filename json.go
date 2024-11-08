package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Respond with 5xx error: ", msg)
	}

	type errResponse struct {
		Error string `json:"error"` // Adds a JSON tag to help with later JSON conversion (no space after :)
	}

	// Respond with an errResponse struct that contains the error message
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload) // Returns the JSON encoding of payload in bytes
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500) // Internal server error
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code) // Write response code
	w.Write(data)       // Write response body
}
