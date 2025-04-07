package main

import (
	"net/http"
	"time"

	"github.com/mcdotjs/chirpy/internal/auth"
)

func (c *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	//NOTE: this is refresh token in Authorization header
	refreshTokenInHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer missing in header", err)
		return
	}
	// refreshToken, err := c.db.GetRefreshTokenByToken(r.Context(), refreshTokenInHeader)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Bearer doesnt exits", err)
	// 	return
	// }
	// if refreshToken.ExpiresAt.Before(time.Now()) {
	// 	//if time.Now().After(refreshToken.ExpiresAt) {
	// 	respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", err)
	// 	return
	// }
	// if refreshToken.RevokedAt.Valid {
	// 	respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", err)
	// 	return
	// }
	user, err := c.db.GetUserFromRefreshToken(r.Context(), refreshTokenInHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}
	accesToken, err := auth.MakeJWT(user.ID, c.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with making new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": accesToken,
	})
}
