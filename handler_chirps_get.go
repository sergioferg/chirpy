package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/sergioferg/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps := []database.Chirp{}

	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorIDUuid, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Malformed user id", err)
			return
		}
		chirps, err = apiCfg.db.GetChirpsFromUser(r.Context(), authorIDUuid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "This user has no chirps", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	} else {
		var err error

		chirps, err = apiCfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
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

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "desc" {
		sort.Slice(chirpSlice, func(i, j int) bool {
			return chirpSlice[i].CreatedAt.After(chirpSlice[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpSlice)
}

func (apiCfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := apiCfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
