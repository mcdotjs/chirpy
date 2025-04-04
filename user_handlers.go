package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	user, err := c.db.CreateUser(r.Context(), params.Email)

	if err != nil {
		log.Fatalf("problem with creating user %s", err)
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time.UTC(),
		UpdatedAt: user.UpdatedAt.Time.UTC(),
	})
}
