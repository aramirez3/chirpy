-- name: CreateUser :one
INSERT INTO users(
    id,
    created_at,
    updated_at,
    email,
    hashed_password
)
    VALUES($1, $2, $3, $4, $5)
    returning *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUsersCount :one
SELECT count(*) FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
    WHERE email=$1;

-- name: GetUserById :one
SELECT * FROM users
    WHERE id=$1;

-- name: UpdateUser :one
UPDATE users
    SET email=$2,
        hashed_password=$3,
        updated_at=$4
    WHERE id = $1
    RETURNING *;

-- name: UpgradeUserToRed :one
UPDATE users
    SET is_chirpy_red=$2,
        updated_at=$3
    WHERE id=$1
    RETURNING *;