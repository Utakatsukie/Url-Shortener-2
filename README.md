# URL Shortener

A backend service for creating and managing short URLs.

The application accepts a long URL, generates a unique 6-character hash, and creates a short URL that can be used to redirect to the original address. Links are stored in PostgreSQL and can be managed through a REST API.

The project is built with Go and follows a layered structure with handlers, services, repositories, middleware, and dependency injection through constructors.

> **Project status:** Backend MVP in active development. Frontend and automated tests are planned.

---

## Features

* Create short URLs
* Redirect from a short URL to the original address
* Unique 6-character URL hashes
* URL validation
* PostgreSQL persistence
* CRUD operations for links
* User registration and authentication
* JWT-based authorization
* Password hashing with bcrypt
* Pagination for retrieving links
* Click statistics
* Middleware-based request processing
* Configuration through environment variables
* PostgreSQL deployment with Docker Compose
* Automatic database migrations with GORM

---

## Tech Stack

| Technology          | Purpose                          |
| ------------------- | -------------------------------- |
| **Go 1.25**         | Backend application              |
| **PostgreSQL 16.4** | Persistent data storage          |
| **GORM**            | ORM and database access          |
| **JWT**             | Authentication and authorization |
| **bcrypt**          | Password hashing                 |
| **Docker Compose**  | PostgreSQL containerization      |
| **net/http**        | HTTP server and routing          |
| **godotenv**        | Environment configuration        |

The project uses `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`, GORM, and the PostgreSQL driver.

---

## How It Works

The main flow of the application is:

```text
Client
   │
   ▼
HTTP Request
   │
   ▼
Middleware
   │
   ├── CORS
   ├── Logging
   └── JWT Authorization
   │
   ▼
HTTP Handler
   │
   ▼
Service / Repository
   │
   ▼
PostgreSQL
```

For a short URL:

```text
Original URL
    │
    ▼
POST /link
    │
    ▼
Generate 6-character hash
    │
    ▼
Store URL + hash in PostgreSQL
    │
    ▼
Short URL
    │
    ▼
GET /{hash}
    │
    ▼
Find original URL
    │
    ├── Record click
    │
    └── HTTP 307 Redirect
    ▼
Original URL
```

The application currently listens on port `8081`.

---

## URL Generation

Each new link receives a randomly generated **6-character hash**.

The hash is generated from Latin letters in both lowercase and uppercase. Before saving a link, the application checks whether the generated hash already exists. If a collision occurs, a new hash is generated until a unique value is found. The hash has a unique database index as an additional constraint.

Example:

```text
Original:
https://example.com/some/very/long/path

Generated hash:
aB7xKp

Short URL:
http://localhost:8081/aB7xKp
```

---

# API

The API is organized into two main areas:

* `/auth` — authentication and registration
* `/link` — link management

The short URL redirect is available directly through `/{hash}`.

## Authentication

### Register

```http
POST /auth/register
```

Request body:

```json
{
  "email": "user@example.com",
  "password": "password",
  "name": "John"
}
```

The password is hashed with bcrypt before being stored in the database.

On successful registration, the API returns a JWT token.

Response:

```json
{
  "token": "<JWT>"
}
```

---

### Login

```http
POST /auth/login
```

Request body:

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

On successful authentication, the API returns a JWT token.

Response:

```json
{
  "token": "<JWT>"
}
```

The JWT contains the authenticated user's email and is signed using the application's configured secret.

---

## Authorization

Protected endpoints require the JWT to be passed using the standard Bearer authentication scheme:

```http
Authorization: Bearer <JWT>
```

JWT validation is implemented as middleware. After successful validation, the authenticated user's email is added to the request context for use by subsequent handlers.

### Currently protected endpoints

| Method  | Endpoint     | JWT      |
| ------- | ------------ | -------- |
| `PATCH` | `/link/{id}` | Required |
| `GET`   | `/link`      | Required |

### Currently public link endpoints

| Method   | Endpoint     | JWT          |
| -------- | ------------ | ------------ |
| `POST`   | `/link`      | Not required |
| `DELETE` | `/link/{id}` | Not required |
| `GET`    | `/{hash}`    | Not required |

> **Note:** Authorization for link creation and deletion is planned for a future version.

---

# Link API

## Create a short link

```http
POST /link
```

Request body:

```json
{
  "url": "https://example.com/very/long/url"
}
```

The URL is validated before being processed.

The server generates a unique 6-character hash and stores the link in PostgreSQL.

Response example:

```json
{
  "id": 1,
  "url": "https://example.com/very/long/url",
  "hash": "aB7xKp"
}
```

---

## Update a link

```http
PATCH /link/{id}
```

**Authentication required.**

Request body:

```json
{
  "url": "https://example.com/new-url",
  "hash": "aB7xKp"
}
```

The endpoint can update the stored URL and hash.

---

## Delete a link

```http
DELETE /link/{id}
```

The endpoint removes the specified link from the database.

**Current status:** authentication is not required yet. Authorization for deletion is planned.

The project uses GORM soft deletion through `gorm.Model`, so deleted records are excluded from normal link queries.

---

## Redirect to the original URL

```http
GET /{hash}
```

Example:

```text
GET /aB7xKp
```

The server:

1. Finds the link by its hash.
2. Records a click.
3. Redirects the client to the original URL.

The redirect uses HTTP **307 Temporary Redirect**.

If the hash does not exist, the API returns `404 Not Found`.

---

## Get links

```http
GET /link?limit=10&offset=0
```

**Authentication required.**

The endpoint supports pagination through `limit` and `offset`.

Response structure:

```json
{
  "links": [],
  "count": 0
}
```

`count` contains the total number of non-deleted links.

> The current implementation returns active links from the database rather than filtering them by the authenticated user. User-specific link ownership and authorization are planned improvements.

---

# Authentication Flow

The authentication flow is intentionally simple:

```text
Register / Login
       │
       ▼
Validate credentials
       │
       ▼
bcrypt password verification
       │
       ▼
Generate JWT
       │
       ▼
Client stores token
       │
       ▼
Authorization: Bearer <JWT>
       │
       ▼
JWT Middleware
       │
       ▼
Validate token
       │
       ▼
Pass request to protected handler
```

Passwords are never stored in plaintext. During registration, the password is hashed using bcrypt with its default cost. During login, the supplied password is compared against the stored hash.

---

# Project Structure

```text
Url-Shortener-2/
│
├── cmd/
│   └── main.go
│
├── configs/
│   └── config.go
│
├── internal/
│   ├── auth/
│   │   ├── errors.go
│   │   ├── handler.go
│   │   ├── payload.go
│   │   └── service.go
│   │
│   ├── link/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── payload.go
│   │   └── repository.go
│   │
│   ├── stat/
│   └── user/
│
├── middleware/
│   ├── auth.go
│   ├── chain.go
│   ├── common.go
│   ├── cors.go
│   └── logs.go
│
├── migrations/
│   └── auto.go
│
├── pkg/
│   ├── db/
│   ├── di/
│   ├── jwt/
│   ├── req/
│   └── res/
│
├── docker-compose.yml
├── go.mod
└── go.sum
```

The project separates application-specific functionality under `internal`, reusable infrastructure under `pkg`, request processing under `middleware`, and database migration logic under `migrations`.

---

# Architecture

The application follows a layered approach.

### Handler

Responsible for:

* receiving HTTP requests;
* parsing request data;
* validating input;
* returning HTTP responses;
* connecting HTTP routes with application logic.

### Service

Contains business logic that does not belong directly in the HTTP layer.

For example, authentication logic is separated into `AuthService`, while the handler remains responsible for HTTP-specific concerns.

### Repository

Responsible for communication with PostgreSQL through GORM.

For example, `LinkRepository` provides operations for:

* creating links;
* finding links by ID;
* finding links by hash;
* updating links;
* deleting links;
* retrieving links with pagination;
* counting active links.

### Middleware

Middleware is used for cross-cutting HTTP functionality:

* JWT authorization;
* CORS;
* request logging;
* middleware composition.

The middleware chain allows these responsibilities to remain separate from individual handlers.

---

# Dependency Injection

The project uses constructor-based dependency injection.

For example, repositories and services are initialized in `main.go` and passed to their respective handlers:

```text
Database
   │
   ├── LinkRepository
   ├── UserRepository
   └── StatRepository
          │
          ▼
     AuthService
          │
          ▼
       Handlers
```

This approach reduces direct dependencies between components and makes the application structure easier to extend and test.

The constructors used throughout the project also make future unit testing more straightforward because dependencies can be replaced with mocks or test implementations.

---

# Configuration

Application configuration is loaded from environment variables.

The configuration layer currently requires:

| Variable | Description                        |
| -------- | ---------------------------------- |
| `DSN`    | PostgreSQL connection string       |
| `TOKEN`  | Secret used for signing JWT tokens |

The application validates that both variables are present during startup.

Create a `.env` file in the project root:

```env
DSN=postgres://postgres:my_password@localhost:5432/postgres?sslmode=disable
TOKEN=your-secret-key
```

---

# Running Locally

## Requirements

Make sure the following are installed:

* Go 1.25+
* Docker
* Docker Compose

---

## 1. Clone the repository

```bash
git clone <repository-url>
cd Url-Shortener-2
```

---

## 2. Start PostgreSQL

The project includes a Docker Compose configuration for PostgreSQL 16.4.

```bash
docker compose up -d
```

The database is exposed on port `5432`.

Check that the container is running:

```bash
docker ps
```

---

## 3. Configure environment variables

Create a `.env` file:

```env
DSN=postgres://postgres:my_password@localhost:5432/postgres?sslmode=disable
TOKEN=your-secret-key
```

The values should match your local PostgreSQL configuration.

---

## 4. Run database migrations

The project currently uses GORM's `AutoMigrate` for creating/updating the database schema.

Run:

```bash
go run migrations/auto.go
```

The migration initializes the tables used by links, users, and click statistics.

---

## 5. Start the application

```bash
go run ./cmd
```

The server will start on:

```text
http://localhost:8081
```

---

# Docker

Docker Compose is currently used to run PostgreSQL separately from the Go application.

```text
┌─────────────────────────┐
│      Go Application     │
│       :8081             │
└────────────┬────────────┘
             │
             │ PostgreSQL
             ▼
┌─────────────────────────┐
│     PostgreSQL 16.4     │
│       Docker            │
│       :5432             │
└─────────────────────────┘
```

The current Compose configuration creates a PostgreSQL container named `postgres_go` and persists database data in `./postgres-data`.

---

# Database

The project uses PostgreSQL as its primary persistent storage and GORM as the ORM.

The main domain entities include:

* **User** — application user and authentication data;
* **Link** — original URL and generated hash;
* **Stat** — click statistics associated with links.

Links use a unique database index for their hash, while GORM's `DeletedAt` field provides soft deletion support.

---

# Current API Reference

| Method   | Endpoint         | Auth | Description                  |
| -------- | ---------------- | ---- | ---------------------------- |
| `POST`   | `/auth/register` | —    | Register a new user          |
| `POST`   | `/auth/login`    | —    | Authenticate and receive JWT |
| `POST`   | `/link`          | —    | Create a short link          |
| `PATCH`  | `/link/{id}`     | JWT  | Update a link                |
| `DELETE` | `/link/{id}`     | —    | Delete a link                |
| `GET`    | `/link`          | JWT  | Get links with pagination    |
| `GET`    | `/{hash}`        | —    | Redirect to original URL     |

---

# Development Status

The project is currently being developed as a backend-focused application.

### Implemented

* [x] REST API
* [x] URL shortening
* [x] PostgreSQL integration
* [x] GORM integration
* [x] User registration
* [x] User authentication
* [x] JWT authorization
* [x] bcrypt password hashing
* [x] Link CRUD operations
* [x] URL validation
* [x] Click statistics
* [x] Middleware chain
* [x] Environment-based configuration
* [x] Dockerized PostgreSQL
* [x] Automatic database migrations

### Planned

* [ ] Automated tests
* [ ] Authorization for link creation
* [ ] Authorization for link deletion
* [ ] User-specific link ownership
* [ ] More granular access control
* [ ] Frontend application
* [ ] Improved migration workflow
* [ ] API documentation with OpenAPI / Swagger
* [ ] Further improvements to URL/hash generation

---

# Why This Project?

This project was created to practice and demonstrate backend development with Go and to explore the design of a small RESTful service.

The main areas of focus are:

* designing REST API endpoints;
* separating HTTP, business, and persistence layers;
* working with PostgreSQL and GORM;
* implementing JWT authentication;
* securely storing passwords with bcrypt;
* creating reusable HTTP middleware;
* using dependency injection and constructor patterns;
* working with Docker-based development environments;
* designing a database-backed URL shortening service.

The project is intentionally being developed incrementally, with authentication, authorization, testing, and infrastructure improvements planned as the application evolves.

---

## License

This project is currently intended primarily as a learning and portfolio project.
