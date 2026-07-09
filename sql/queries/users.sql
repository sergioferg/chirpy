-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC',
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;
--

-- name: Reset :exec
DELETE FROM users;
--

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
--

-- name: UpdateUser :one
UPDATE users
SET email = $1,
    hashed_password = $2,
    updated_at = NOW() AT TIME ZONE 'UTC'
WHERE id = $3
RETURNING id, email, created_at, updated_at;
--
