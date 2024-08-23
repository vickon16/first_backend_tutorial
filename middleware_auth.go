package main

import (
	"fmt"
	"net/http"

	"github.com/vickon16/first-go-backend/internal/auth"
	"github.com/vickon16/first-go-backend/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *apiConfig) middleWareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Api Key Error: %s\n", err.Error()))
			return
		}

		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %s\n", err.Error()))
			return
		}

		handler(w, r, user)
	}
}
