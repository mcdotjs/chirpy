package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mcdotjs/chirpy/internal/database"
)

func (c *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {

	type CreateChirpParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := CreateChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	cleanedBody := checkForProfane(params.Body)

	newChirp := &database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	}
	chirp, err := c.db.CreateChirp(r.Context(), *newChirp)

	if err != nil {
		log.Fatalf("problem with creating chirp %s", err)
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.UTC(),
		UpdatedAt: chirp.UpdatedAt.UTC(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
