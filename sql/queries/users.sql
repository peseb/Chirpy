-- name: CreateUser :one
INSERT INTO users (id, created_at, modified_at, email)
VALUES (
    gen_rangom_uuid(), NOW(), NOW(), $1
)
RETURNING *;