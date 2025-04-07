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

type returnedUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"token,omitempty"`
}

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type loginParameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	expiresInSeconds := 3600
	decoder := json.NewDecoder(r.Body)
	params := loginParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "LOGIN: Error decoding parameters", err)
		return
	}

	if params.ExpiresInSeconds > 0 {
		if params.ExpiresInSeconds > 3600 {
			expiresInSeconds = 3600
		} else {
			expiresInSeconds = params.ExpiresInSeconds
		}
	}

	userByEmail, err := c.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password 1", err)
		return
	}

	res := auth.CheckPasswordHash(params.Password, userByEmail.HashedPassword)
	if res != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password 2", err)
		return
	}

	generatedToken, err := auth.MakeJWT(userByEmail.ID, c.jwtSecret, time.Duration(expiresInSeconds)*time.Second)
	fmt.Println("generatedToken", generatedToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with jwt generation", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnedUser{
		ID:        userByEmail.ID,
		Email:     userByEmail.Email,
		CreatedAt: userByEmail.CreatedAt.UTC(),
		UpdatedAt: userByEmail.UpdatedAt.UTC(),
		Token:     generatedToken,
	})
}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hashedPassoword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("hashing error", err)
	}

	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassoword,
	})

	if err != nil {
		log.Fatalf("problem with creating user %s", err)
	}

	respondWithJSON(w, http.StatusCreated, returnedUser{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.UTC(),
		UpdatedAt: user.UpdatedAt.UTC(),
	})
}
