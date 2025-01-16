package main

import (
	"net/http"

	"github.com/FedjaW/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
    chirpIDString := r.PathValue("chirpID")
    id, err := uuid.Parse(chirpIDString)
    if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

    token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

    chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

    if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Couldn't get chirp", err)
		return
    }

    cfg.db.DeleteChirp(r.Context(), id)

	respondWithJSON(w, http.StatusNoContent, nil)
}
