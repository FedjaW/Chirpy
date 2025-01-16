package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
    type data struct {
        UserId uuid.UUID `json:"user_id"`
    }
    type parameters struct {
        Event string `json:"event"`
        Data data `json:"data"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
    }

    if params.Event != "user.upgraded" {
        respondWithJSON(w, http.StatusNoContent, nil)
        return
    }

    _, err = cfg.db.MakeUserRed(r.Context(), params.Data.UserId)
    if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user", err)
		return
    }

    respondWithJSON(w, http.StatusNoContent, nil)
}
