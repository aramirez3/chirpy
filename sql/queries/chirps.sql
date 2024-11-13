-- name: CreateChirp :one

INSERT INTO chirps(
    id,
    created_at,
    updated_at,
    body,
    user_id
)
    VALUES($1, $2, $3, $4, $5)
    RETURNING *;

-- name: DeleteAllChirps :exec

DELETE FROM chirps;

-- name: GetChirpById :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirpsCount :one
SELECT count(*) FROM chirps;

-- name: GetAllChirpsAsc :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirpByAuthorIdAsc :many
SELECT * FROM chirps
    WHERE user_id=$1
    ORDER BY created_at ASC;

-- name: GetChirpByAuthorIdDesc :many
SELECT * FROM chirps
    WHERE user_id=$1
    ORDER BY created_at DESC;

-- name: DeleteChirpById :one
DELETE FROM chirps
    WHERE id=$1
    RETURNING *;