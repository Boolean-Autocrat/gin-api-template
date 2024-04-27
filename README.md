# Gin API Template (Go)

This is a template/boilterplate for a REST API using Gin, Postgres, SQLC and Docker built with Go (meant to be used as a starting point for new projects).

## Requirements

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [SQLC](https://sqlc.dev/)
- [Make](https://www.gnu.org/software/make/)
- [Air](https://github.com/cosmtrek/air)
- [Golang Migrate](https://github.com/golang-migrate/migrate)

## Features

- Basic authentication using Google OAuth2 (with classic cookie session)
- Pre built middlewares for CORS and authentication
- SQLC configured for type-safe SQL queries
- Dockerized development environment
- Fully dockerized production environment with multi-stage build for minimal image size, pgadmin for database management and managed migrations
- Clean file structure for the project

## Getting Started

1. Clone the repository

```bash
git clone https://github.com/Boolean-Autocrat/gin-api-template.git <project_name>
cd <project_name>
```

3. Copy the `.env.example` file to `.env` and update the values.

```bash
cp .env.example .env
```

4. Spawn the postgres development database

```bash
make devdb
```

5. Update the database url in the `Makefile` as per your configuration.

6. Run the migrations

```bash
make migrateup
```

7. Run the server

```bash
air
```

8. Visit `http://localhost:3000/health` or `http://localhost:3000/example` in your browser to see the application.

## Makefile Commands

- `make devdb` - Start the development database
- `make devdbdown` - Stop the development database
- `make migratecreate name=<migration_name>` - Create a new migration
- `make migrateup` - Run the migrations
- `make migratedown` - Rollback the migrations

## Included Endpoints

- `health` - Health check endpoint
- `example` - Example endpoint
- `auth/google/login` - Google OAuth2 login endpoint
- `auth/google/callback` - Google OAuth2 callback endpoint
- `auth/google/logout` - Google OAuth2 logout endpoint

## Environment Variables

- `SESSION_SECRET` - Secret key for session management
- `POSTGRES_HOST` - Hostname for the postgres database
- `POSTGRES_PORT` - Port for the postgres database
- `POSTGRES_USER` - Username for the postgres database
- `POSTGRES_PASSWORD` - Password for the postgres database
- `POSTGRES_DB` - Database name for the postgres database
- `DB_SOURCE` - Connection string for the postgres database
- `PORT` - Port for the gin server
- `GIN_MODE` - Mode for the gin server (debug/release)
- `PGADMIN_DEFAULT_EMAIL` - Default email for pgadmin
- `PGADMIN_DEFAULT_PASSWORD` - Default password for pgadmin
- `GOOGLE_CLIENT_ID` - Google OAuth2 client id
- `GOOGLE_CLIENT_SECRET` - Google OAuth2 client secret
- `GOOGLE_REDIRECT_URL` - Google OAuth2 redirect url
- `LOGIN_REDIRECT_URL` - Redirect url after login
- `LOGOUT_REDIRECT_URL` - Redirect url after logout
- `COOKIE_SET_URL` - Domain to set the cookie

## File Structure

```
│   .air.toml
│   .dockerignore
│   .env.example
│   .gitignore
│   docker-compose.dev.yml
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│   main.go
│   Makefile
│   README.md
│   sqlc.yaml
│   start.sh
│   tree.txt
│   wait-for-it.sh
│
├───api
│   ├───auth
│   │       auth.go
│   │
│   ├───example
│   │       example.go
│   │
│   └───utils
│           utils.go
│
├───db
│   ├───migrations
│   │       000001_init.down.sql
│   │       000001_init.up.sql
│   │
│   ├───query
│   │       auth.sql
│   │
│   └───sqlc
│           auth.sql.go
│           db.go
│           models.go
│           postgres.go
│
├───middlewares
│       authMiddleware.go
│       corsMiddleware.go
│       verifySession.go
```
