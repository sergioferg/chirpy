package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	type validResp struct {
		CleanBody string `json:"cleaned_body"`
	}

	respValid := validResp{
		CleanBody: censorChirp(chirp.Body),
	}

	respondWithJSON(w, http.StatusOK, respValid)
}

func censorChirp(str string) string {
	parsed := strings.Fields(str)
	for i, word := range parsed {
		if checkProfanity(word) {
			parsed[i] = "****"
		}
	}

	return strings.Join(parsed, " ")
}

func checkProfanity(word string) bool {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	return slices.Contains(profaneWords, strings.ToLower(word))
}
