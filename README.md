# Go URL Saver API

A backend CRUD API built with Go and PostgreSQL using a layered architecture.

This project started as a simple in-memory URL saver and was later upgraded step-by-step into a PostgreSQL-backed API with Docker support. The main focus was to understand backend fundamentals properly instead of just following tutorials.

---

## Features

- Create, read, update, and delete saved links
- PostgreSQL persistence
- Docker and Docker Compose setup
- Environment variable configuration
- Middleware logging
- Graceful shutdown
- Layered architecture (handler → service → store)

---

## Tech Stack

- **Go** (standard library `net/http`, `database/sql`)
- **PostgreSQL** (via `lib/pq` driver)
- **Docker** and **Docker Compose**

---

## Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go
└── internal/
    ├── config/
    ├── handlers/
    ├── middleware/
    ├── models/
    ├── routes/
    ├── service/
    └── store/
```

---

## Architecture

The project follows a layered structure:

```
Handler → Service → Store → PostgreSQL
```

**Handler Layer** — HTTP concerns only. Parses requests, encodes/decodes JSON, sends responses.

**Service Layer** — Business logic. Validation, URL formatting, and other rules live here.

**Store Layer** — Database access. Handles SQL queries and persistence details.

The app was initially built with an in-memory store, then swapped to PostgreSQL without changing the handler or service layers. That refactor proved the architecture was working.

---

## API Endpoints

### Create a link

```
POST /save_url
```

**Request Body:**

```json
{
  "title": "Google",
  "link": "https://google.com"
}
```

### Get all links

```
GET /get_all
```

### Update a link

```
PUT /update_link?id=1
```

**Request Body:**

```json
{
  "title": "Updated Title",
  "link": "https://example.com"
}
```

### Delete a link

```
DELETE /delete_link?id=1
```

---

## Environment Variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
PORT=8080
```

---

## Running Locally

Make sure PostgreSQL is installed and running locally, then:

```bash
go run ./cmd/api
```

---

## Running with Docker

Build and start both the Go API and PostgreSQL containers:

```bash
docker compose up --build
```

To stop everything:

```bash
docker compose down
```

---

## What I Learned

- Layered backend architecture in Go
- Dependency injection by hand (no frameworks)
- PostgreSQL CRUD operations (`QueryRow`, `Query`, `Exec`)
- Structuring Go projects with `internal/` packages
- Error handling patterns without exceptions
- Docker networking between services
- Environment configuration with `.env`
- Graceful shutdown with signal handling

---

## Future Improvements

- Unit tests for the service layer
- JWT-based authentication
- Pagination for the `get_all` endpoint
- Redis caching layer

---

## Why I Built This

I built this project to strengthen my backend fundamentals before moving deeper into Java and Spring Boot.

I wanted to understand:

- How APIs actually work under the hood
- How databases connect to real applications
- Why layering and abstractions matter
- How containers and environment configs fit together

I plan to rebuild the same project in Java/Spring Boot later to compare both ecosystems and the tradeoffs each one makes.