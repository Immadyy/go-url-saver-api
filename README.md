# Go URL Saver API

A backend CRUD API built with Go and PostgreSQL using a layered architecture.

This project started as a simple in-memory URL saver and was later upgraded step-by-step into a PostgreSQL-backed API with Docker support. The main focus of the project was to understand backend fundamentals properly instead of only following tutorials.

---

# Features

- Create links
- Get all saved links
- Update existing links
- Delete links
- PostgreSQL integration
- Docker + Docker Compose setup
- Environment variable configuration
- Middleware logging
- Graceful shutdown
- Layered architecture

---

# Tech Stack

- Go
- PostgreSQL
- Docker
- Docker Compose
- net/http
- database/sql
- lib/pq

---

# Project Structure

```text
cmd/
    api/
        main.go

internal/
    config/
    handlers/
    middleware/
    models/
    routes/
    service/
    store/

-------------------------------------------------------------------------------------------------------------

    Architecture:

The project follows a layered structure:

Handler -> Service -> Store -> PostgreSQL
Handler Layer

Responsible for:

HTTP requests
decoding/encoding JSON
sending responses
Service Layer

Responsible for:

business logic
validation
URL formatting/checking
Store Layer

Responsible for:

SQL queries
database interaction
persistence

The project was first built using an in-memory store and later switched to PostgreSQL without changing the whole application structure.

API Endpoints
Create Link
POST /save_url

Request Body:

{
  "title": "Google",
  "link": "https://google.com"
}
Get All Links
GET /get_all
Update Link
PUT /update_link?id=1

Request Body:

{
  "title": "Updated Title",
  "link": "https://example.com"
}
Delete Link
DELETE /delete_link?id=1
Environment Variables

Create a .env file in the project root.

DB_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
PORT=8080
Run Without Docker

Make sure PostgreSQL is installed and running locally.

Run the server:

go run ./cmd/api
Run With Docker

Build and start the containers:

docker compose up --build

This starts:

Go API container
PostgreSQL container

Stop containers:

docker compose down
What I Learned From This Project
Layered backend architecture
Dependency injection in Go
PostgreSQL CRUD operations
QueryRow vs Query vs Exec
Error handling patterns
Docker networking basics
Environment configuration
Graceful shutdown handling
Future Improvements
Unit tests for service layer
JWT authentication
Pagination
Redis caching
Why I Built This

I built this project to strengthen my backend fundamentals before moving deeper into Java and Spring Boot.

The goal was to understand:

how APIs work internally
how databases connect to applications
how layering and abstractions work
how containers and environments work

I also plan to rebuild the same project later in Java/Spring Boot to compare both ecosystems and design tradeoffs.