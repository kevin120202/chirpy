# Chirpy API

A modern, RESTful API for a social media platform built with Go, featuring user authentication, chirp (post) management, and webhook integrations.

## 🚀 Features

- **User Management**: Registration, authentication, and profile updates
- **Chirp System**: Create, read, and delete social media posts
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Content Moderation**: Automatic profanity filtering for chirps
- **Webhook Integration**: Support for external service integrations
- **Premium Features**: Chirpy Red subscription system
- **Database**: PostgreSQL with SQLC for type-safe database operations

## 🛠️ Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **ORM**: SQLC (SQL Compiler)
- **Authentication**: JWT with bcrypt password hashing
- **HTTP Server**: Standard library `net/http`

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE
);
```

### Chirps Table
```sql
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);
```

## 🔌 API Endpoints

### Authentication Endpoints

#### POST /api/users
Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false
  }
}
```

#### POST /api/login
Authenticate a user and receive access tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false
  },
  "token": "jwt-access-token",
  "refresh_token": "refresh-token"
}
```

#### POST /api/refresh
Refresh an access token using a refresh token.

**Request Body:**
```json
{
  "refresh_token": "refresh-token"
}
```

**Response:**
```json
{
  "token": "new-jwt-access-token"
}
```

#### POST /api/revoke
Revoke a refresh token.

**Request Body:**
```json
{
  "refresh_token": "refresh-token"
}
```

### Chirp Endpoints

#### POST /api/chirps
Create a new chirp (requires authentication).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "body": "Hello, world! This is my first chirp!"
}
```

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello, world! This is my first chirp!",
  "user_id": "user-uuid"
}
```

#### GET /api/chirps
Get all chirps (optional query parameters).

**Query Parameters:**
- `author_id`: Filter by user ID
- `sort`: Sort order ("asc" or "desc", default: "desc")

**Response:**
```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "Hello, world!",
    "user_id": "user-uuid"
  }
]
```

#### GET /api/chirps/{chirpID}
Get a specific chirp by ID.

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello, world!",
  "user_id": "user-uuid"
}
```

#### DELETE /api/chirps/{chirpID}
Delete a chirp (requires authentication and ownership).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

### User Management Endpoints

#### PUT /api/users
Update user information (requires authentication).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

### Webhook Endpoints

#### POST /api/polka/webhooks
Handle webhook events (requires API key authentication).

**Headers:**
```
Authorization: ApiKey <api-key>
```

**Request Body:**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "user-uuid"
  }
}
```

### System Endpoints

#### GET /api/healthz
Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

#### GET /admin/metrics
Get server metrics (requires no authentication).

**Response:**
```json
{
  "hits": 42
}
```

#### POST /admin/reset
Reset server metrics (requires no authentication).

## 🔐 Authentication

### JWT Tokens
The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <jwt-token>
```

### Refresh Tokens
Refresh tokens are used to obtain new access tokens without re-authentication. They have a longer lifespan and can be revoked.

### API Keys
For webhook endpoints, use API key authentication:

```
Authorization: ApiKey <api-key>
```

## 📊 Chirpy Red Features

Users can be upgraded to "Chirpy Red" status through webhook events, which may unlock premium features in the future.

## 🏗️ Project Structure

```
chirpy/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── sqlc.yaml              # SQLC configuration
├── .env                   # Environment variables (create this)
├── internal/
│   ├── auth/              # Authentication utilities
│   │   ├── auth.go        # JWT and password functions
│   │   └── auth_test.go   # Authentication tests
│   └── database/          # Generated database code
│       ├── db.go          # Database connection
│       ├── models.go      # Database models
│       └── *.sql.go       # Generated query functions
├── sql/
│   ├── schema/            # Database migrations
│   │   ├── 001_users.sql
│   │   ├── 002_chirps.sql
│   │   ├── 003_add_hashed_password.sql
│   │   ├── 004_refresh_tokens.sql
│   │   └── 005_chirpy_red.sql
│   └── queries/           # SQL queries for SQLC
│       ├── users.sql
│       ├── chirps.sql
│       ├── refresh_token.sql
│       └── reset.sql
├── handler_*.go           # HTTP request handlers
├── metrics.go             # Metrics middleware
├── json.go                # JSON response utilities
├── readiness.go           # Health check handler
├── reset.go               # Metrics reset handler
└── assets/                # Static assets
```
