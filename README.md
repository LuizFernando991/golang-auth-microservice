# Auth Microservice (Go + Gin + Postgres + Redis)

A sample authentication microservice built with Go using Gin, PostgreSQL, and Redis (for rate limiting).
Designed with best practices in mind: layered architecture, dependency injection, context usage, JWT with rotating refresh tokens, bcrypt password hashing, database migrations, Docker support, and testing.

## Features
- User registration and login with secure password hashing (bcrypt)

- JWT-based authentication with refresh tokens

- Rate limiting on authentication endpoints using Redis

- PostgreSQL integration with connection pooling

- Clean layered architecture (handlers, services, repositories)

- Configuration via environment variables

- Dockerized setup for easy deployment and development

- Automated database migrations

- Unit and integration tests

## Getting Started
### Prerequisites
- Docker & Docker Compose

- Go 1.20+ (for local builds if not using Docker)

- PostgreSQL (if not using Docker)

- Redis (if not using Docker)

## Installation
1. Clone the repository
 ```bash
git clone https://github.com/LuizFernando991/golang-auth-microservice.git
cd golang-auth-microservice
```
2. Configure environment variables
```bash
cp .env.example .env
```
3. Start with Docker Compose

```bash 
docker compose up --build
```

## Running Database Migrations
To run the database migrations inside the running PostgreSQL container (if you're not using the docker compose):

```bash
docker exec -i <POSTGRES_CONTAINER_NAME> psql -U postgres -d authdb < internal/db/migrations.sql
```
Replace <POSTGRES_CONTAINER_NAME> with the actual container name (e.g., auth_postgres_1).

## Rate Limiting
Rate limiting is applied on all auth endpoints (/register, /login, /refresh, /logout) using Redis to prevent abuse.

## Project Structure

```bash
.
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── config
│   ├── db
│   │   └── migrations.sql
│   ├── handler
│   ├── middleware
│   ├── model
│   ├── repository
│   ├── server
│   ├── service
│   └── util
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── README.md
```

## Examples of Requests

1. Register a New User - POST /v1/register

```bash
curl -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{
        "email": "user@example.com",
        "password": "strongpassword123"
      }'
```

2. Login User - POST /v1/login

```bash
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{
        "email": "user@example.com",
        "password": "strongpassword123"
      }'
```

3. Refresh Tokens

```bash
curl -X POST http://localhost:8080/v1/refresh \
  -H "Refresh_Token: Bearer your_refresh_token_here"
```

4. Logout User

```bash
curl -X POST http://localhost:8080/v1/logout \
  -H "Refresh_Token: Bearer your_access_token_here"
```

5. Get Current User Info - GET /v1/me

```bash
curl -X GET http://localhost:8080/v1/me \
  -H "Authorization: Bearer your_access_token_here"
```

## Testing

```bash
go test ./internal/service/...
```


