package main

import (
	"net/http"
)

func (apiCfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirpSlice := make([]Chirp, 0, len(chirps))
	for _, chirp := range chirps {
		chirpSlice = append(chirpSlice, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpSlice)
}
