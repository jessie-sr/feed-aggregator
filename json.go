package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload) // Returns the JSON encoding of payload in bytes
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500) // Internal server error
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200) // Success
	w.Write(data)      // Write response body
}
