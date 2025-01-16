package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FedjaW/Chirpy/internal/auth"
	"github.com/FedjaW/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
        Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
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

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

    createChirpParams := database.CreateChirpParams{
        Body: params.Body,
        UserID: userID,
    }
    chirp, err := cfg.db.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
            ID: chirp.ID,            
            CreatedAt: chirp.CreatedAt,
            UpdatedAt: chirp.UpdatedAt,
            Body: chirp.Body,
            UserID: chirp.UserID,
        },
	})
}
