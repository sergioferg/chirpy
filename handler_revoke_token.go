package main

import (
	"net/http"

	"github.com/sergioferg/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing bearer token", err)
		return
	}

	err = apiCfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid/Missing token", err)
		return
	}

	type response struct{}

	respondWithJSON(w, http.StatusNoContent, response{})
}
