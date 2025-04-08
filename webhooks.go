package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (c *apiConfig) redChirpWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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
