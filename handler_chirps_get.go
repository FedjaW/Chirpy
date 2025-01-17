package main

import (
	"net/http"
	"sort"

	"github.com/FedjaW/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r* http.Request) {
    a := r.URL.Query().Get("author_id")
    var dbChirps []database.Chirp
    var err error
    if a != "" {
        authorID, _ := uuid.Parse(a)
        dbChirps, err = cfg.db.ListChirpsByUser(r.Context(), authorID)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
            return
        }
    } else {
        dbChirps, err = cfg.db.ListChirps(r.Context())
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
            return
        }
    }

    chirps := []Chirp{}
    for _, dbChirp := range dbChirps {
        chirps = append(chirps, Chirp{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            UserID:    dbChirp.UserID,
            Body:      dbChirp.Body,
        })
    }

    sorting := r.URL.Query().Get("sort")
    if sorting == "asc" || sorting == "" {
        sort.Slice(
            chirps, 
            func(i, j int) bool { return chirps[i].CreatedAt.String() < chirps[j].CreatedAt.String() })
    } else {
        sort.Slice(
            chirps, 
            func(i, j int) bool { return chirps[i].CreatedAt.String() > chirps[j].CreatedAt.String() })
    }

    respondWithJSON(w, http.StatusOK, chirps)
}
