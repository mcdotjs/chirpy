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
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	AccessToken  string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type loginParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := loginParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "LOGIN: Error decoding parameters", err)
		return
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

	generatedAccessToken, err := auth.MakeJWT(userByEmail.ID, c.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with jwt generation", err)
		return
	}

	generatedRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with making refresh token", err)
		return
	}

	storedRefreshToken, err := c.db.CreateRefreshToken(r.Context(),
		database.CreateRefreshTokenParams{
			Token:  generatedRefreshToken,
			UserID: userByEmail.ID,
		},
	)

	//NOTE: egal which is in response??
	fmt.Println("stored and generated has to be same", storedRefreshToken.Token, generatedRefreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with storing refresh token in database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnedUser{
		ID:           userByEmail.ID,
		Email:        userByEmail.Email,
		CreatedAt:    userByEmail.CreatedAt.UTC(),
		UpdatedAt:    userByEmail.UpdatedAt.UTC(),
		AccessToken:  generatedAccessToken,
		RefreshToken: generatedRefreshToken,
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
