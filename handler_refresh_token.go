package main

import (
	"net/http"
	"time"

	"github.com/sergioferg/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing bearer token", err)
		return
	}

	user, err := apiCfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid/Expired token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, apiCfg.secret, time.Duration(1*time.Hour))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating JWT token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
