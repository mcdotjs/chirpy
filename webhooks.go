package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mcdotjs/chirpy/internal/auth"
)

func (c *apiConfig) redChirpWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer missing when updating user", err)
		return
	}
	if apiKey != c.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "Bearer missing when updating user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = c.db.UpdateUserToRedChirp(r.Context(), uuid.MustParse(params.Data.UserID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong with updeting user red chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
