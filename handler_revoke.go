package main

import (
	"net/http"

	"github.com/FedjaW/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
    token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token", err)
		return
	}

    _, err = cfg.db.RevokeRefreshToken(r.Context(), token)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
        return
    }
    respondWithJSON(w, http.StatusNoContent, nil)
}
