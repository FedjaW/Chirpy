package main

import (
	"net/http"
	"time"

	"github.com/FedjaW/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
    type response struct {
        Token string `json:"token"`
    }

    token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no token", err)
		return
	}

    user, err:= cfg.db.GetUserFromRefreshToken(r.Context(), token)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't get user", err)
        return
    }

	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't create access JWT", err)
		return
	}

    respondWithJSON(w, http.StatusOK, response{
        Token: accessToken,
    })
}

