-- name: CreateUser :one
INSERT INTO users(
    id,
    created_at,
    updated_at,
    email
)
    values($1, $2, $3, $4)
    returning *;