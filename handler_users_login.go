package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sergioferg/chirpy/internal/auth"
	"github.com/sergioferg/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := apiCfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking hash", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}
	tokenDuration := time.Duration(1 * time.Hour)

	accessToken, err := auth.MakeJWT(user.ID, apiCfg.secret, tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating JWT token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = apiCfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving refresh token to database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
