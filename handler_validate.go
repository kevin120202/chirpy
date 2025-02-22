package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the parameters", err)
		return
    }

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirps is too long", nil)
		return
	}

	cleaned := getCleanedBody(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleaned})
}

func getCleanedBody(chirp string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(chirp, " ")

	for i, w := range words {
		loweredWord := strings.ToLower(w)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	
	cleaned := strings.Join(words, " ")
	return cleaned
}