package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mcdotjs/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem wihth getting all chirps", err)
	}

	finalChirps := make([]Chirp, 0, len(chirps))
	for _, ch := range chirps {
		finalChirps = append(finalChirps, Chirp{
			ID:        ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body:      ch.Body,
			UserID:    ch.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, finalChirps)
}

func (c *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {

	type CreateChirpParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.UTC(),
		UpdatedAt: chirp.UpdatedAt.UTC(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
