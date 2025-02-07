// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteAllUser = `-- name: DeleteAllUser :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUser(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUser)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const makeUserRed = `-- name: MakeUserRed :one
UPDATE users SET is_chirpy_red = true
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

func (q *Queries) MakeUserRed(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, makeUserRed, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserData = `-- name: UpdateUserData :one
UPDATE users SET email = $2,
hashed_password = $3
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateUserDataParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

func (q *Queries) UpdateUserData(ctx context.Context, arg UpdateUserDataParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserData, arg.ID, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
