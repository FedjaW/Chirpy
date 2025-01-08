package main

import (
	"encoding/json"
	"net/http"
    "github.com/google/uuid"
    "time"
    "github.com/FedjaW/Chirpy/internal/database"
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
		UserID uuid.UUID `json:"user_id"`
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

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

    createChirpParams := database.CreateChirpParams{
        Body: params.Body,
        UserID: params.UserID,
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
