package main

import (
	"net/http"

	"github.com/mcdotjs/chirpy/internal/auth"
)

func (c *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	//NOTE: this is refresh token in Authorization header
	refreshTokenInHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer missing in header", err)
		return
	}
	_, err = c.db.GetRefreshTokenByToken(r.Context(), refreshTokenInHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer doesnt exits", err)
		return
	}

	_, err = c.db.UpdateRefreshTokenByToken(r.Context(), refreshTokenInHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
