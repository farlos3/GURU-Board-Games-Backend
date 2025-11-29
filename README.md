# GURU Board Games Backend

This repository contains the backend services for the GURU Board Games platform, including a Go-based API gateway and a Python microservice for recommendations.

## Project Structure

```
GO-Gateway/      # Go API Gateway and main backend logic
Pyservice/       # Python microservice for recommendations
```

### GO-Gateway

- **Language:** Go
- **Main Entry:** `main.go`
- **Features:**
  - User authentication (register, login, OTP, profile)
  - Board game management (CRUD, game rules)
  - Game search and state management
  - Recommendation client integration
  - Database connection and repositories

#### Key Directories

- `internal/auth/` - Auth handlers, JWT, OTP, and service logic
- `internal/boardgame/` - Board game handlers and services
- `internal/db/` - Database connection and repositories
- `internal/gamesearch/` - Game search handlers
- `internal/gamestate/` - Game state handlers
- `internal/recommendation/` - Recommendation client and handler
- `internal/useractivity/` - User activity handlers
- `models/` - Data models
- `routes/` - API route definitions

### Pyservice

- **Language:** Python
- **Main Entry:** `main.py`
- **Features:**
  - Recommendation engine (popularity, mapping, indexing)
  - Database connection
  - Service logic for board game recommendations

#### Key Directories

- `connection/` - Database connection logic
- `recomendation/` - Recommendation algorithms and services

## Getting Started

### Prerequisites

- Go 1.18+
- Python 3.8+
- Docker (optional, for containerization)
- PostgreSQL or compatible database

### Setup

#### GO-Gateway

1. Install dependencies:
   ```bash
   cd GO-Gateway
   go mod tidy
   ```
2. Configure environment variables in `.env`.
3. Run the server:
   ```bash
   go run main.go
   ```
   Or use [Air](https://github.com/cosmtrek/air) for live reload:
   ```bash
   air
   ```

#### Pyservice

1. Install dependencies:
   ```bash
   cd Pyservice
   pip install -r requirement.txt
   ```
2. Configure environment variables in `.env`.
3. Run the service:
   ```bash
   python main.py
   ```

## API Endpoints

- See `GO-Gateway/routes/routes.go` for all available endpoints.
- Auth, board game, game state, recommendation, and user activity APIs are provided.

## Environment Variables

Both services use a `.env` file for configuration. Example variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASS=yourpass
DB_NAME=yourdb
SECRET_KEY=yourkey
```
