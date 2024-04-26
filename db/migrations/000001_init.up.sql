CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "picture" varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "user_sessions" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user" uuid NOT NULL UNIQUE,
  "token" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);