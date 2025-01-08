package main

import (
	"net/http"
    "github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r* http.Request) {
    chirpIDString := r.PathValue("chirpID")
    id, err := uuid.Parse(chirpIDString)
    if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

    type response struct {
        Chirp
    }

    dbChirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
    }

    respondWithJSON(w, http.StatusOK, Chirp{
            ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
        },
    )
}
