package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sergioferg/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	strChirpID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(strChirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid/Missing Token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Expired/Invalid Token", err)
		return
	}

	chirp, err := apiCfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error retrieving chirp", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Invalid token", nil)
		return
	}

	err = apiCfg.db.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
