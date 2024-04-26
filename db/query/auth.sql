-- name: CreateUser :one
INSERT INTO "users" ("name", "email", "picture") VALUES ($1, $2, $3) RETURNING "id";

-- name: GetUserByEmail :one
SELECT * FROM "users" WHERE "email" = $1;

-- name: CreateOrUpdateSession :exec
INSERT INTO "user_sessions" ("user", "token") VALUES ($1, $2) ON CONFLICT ("user") DO UPDATE SET "token" = $2, "created_at" = now();

-- name: GetSession :one
SELECT * FROM "user_sessions" WHERE "token" = $1;

-- name: DeleteSession :exec
DELETE FROM "user_sessions" WHERE "user" = $1 AND "token" = $2;