package main

import (
	"fmt"
	"net/http"

	"github.com/jessie-sr/rss-aggregator/internal/auth"
	"github.com/jessie-sr/rss-aggregator/internal/db"
)

// Custom handler type with an authenticated user as the 3rd param
type authedHandler func(http.ResponseWriter, *http.Request, db.User)

// middlewareAuth wraps an authedHandler to make it compatible with the standard http.HandlerFunc type
func (apiCig *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get API key from the request header
		apikey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Authentication error: %v", err))
			return
		}

		// Get user with the API key
		user, err := apiCig.DB.GetUserByApiKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error getting user: %s", err))
			return
		}

		handler(w, r, user)
	}
}
