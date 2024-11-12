-- name: CreateChirp :one

INSERT INTO chirps(
    id,
    created_at,
    updated_at,
    body,
    user_id
)
    values($1, $2, $3, $4, $5)
    RETURNING *;

-- name: DeleteAllChirps :exec

DELETE from chirps;

-- name: GetChirpsCount :one
Select count(*) from chirps;