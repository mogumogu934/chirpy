-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET updated_at = NOW(), email = $1, hashed_password = $2
WHERE id = $3;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = TRUE
where id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;