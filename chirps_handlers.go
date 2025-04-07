package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mcdotjs/chirpy/internal/auth"
	"github.com/mcdotjs/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("id")
	chirp, err := c.db.GetChirpById(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Problem wihth getting chirp by ID", err)
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
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
		Body string `json:"body"`
	}

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer missing", err)
		return
	}

	id, err := auth.ValidateJWT(bearer, c.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating bearer", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := CreateChirpParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	cleanedBody := checkForProfane(params.Body)

	newChirp := &database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: id,
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

func (c *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	refreshTokenInHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer missing when updating user", err)
		return
	}
	ch, err := c.db.GetChirpById(r.Context(), uuid.MustParse(chirpID))
	fmt.Println("PPPPP", ch)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp was not found", err)
		return
	}

	userId, err := auth.ValidateJWT(refreshTokenInHeader, c.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	_, err = c.db.DeleteChirpById(r.Context(), database.DeleteChirpByIdParams{
		ID:     uuid.MustParse(chirpID),
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusForbidden, "You have no permission for this", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
