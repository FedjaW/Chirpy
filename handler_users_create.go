package main

import (
    "time"
    "github.com/google/uuid"
    "encoding/json"
	"net/http"
    "github.com/FedjaW/Chirpy/internal/database"
    "github.com/FedjaW/Chirpy/internal/auth"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool  `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
        User
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }
    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
        return
    }
    createUserParams := database.CreateUserParams{
        Email: params.Email,
        HashedPassword: hashedPassword,
    }
    user, err := cfg.db.CreateUser(r.Context(), createUserParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
        return
    }

    respondWithJSON(w, http.StatusCreated, response{
        User: User{
            ID: user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email: user.Email,
            IsChirpyRed: user.IsChirpyRed,
        },
    })
}

