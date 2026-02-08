# Ping Pong Backend API Documentation

## 📋 Overview

This is a comprehensive backend server for the Hanoi Super League Ping Pong Tournament system, built with:

- **Go**
- **PostgreSQL** (database)
- **Docker** (containerization)

## 🏗️ Architecture

### Project Structure

```markdown
backend-ping-pong-app/
├── main.go              # Application entry point
├── models.go            # Data models and request/response structs
├── database.go          # Database queries and operations
├── routes.go            # API handlers (renamed from handlers.go)
├── handlers.go          # File upload handlers (legacy)
├── init.sql             # Database initialization script
├── docker-compose.yml   # Docker compose configuration
├── Dockerfile           # Docker build configuration
├── .env.example         # Environment variables template
├── go.mod               # Go module dependencies
└── README.md            # This file
```

### Data Model Layers

1. **Core Tables**: ranks, teams, players, seasons
2. **Season Tables**: player_seasons, player_point_logs
3. **Competition Tables**: fixtures, matches
4. **Import Tables**: staging_players (for bulk imports)

---

## 🚀 Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- PostgreSQL 15+ (if running without Docker)

### Quick Start with Docker

#### **Clone the repository**

```bash
cd backend-ping-pong-app
```

#### **Copy environment file**

```bash
cp .env.example .env
# Modify .env if needed for your environment
```

#### **Build and start services**

```bash
docker-compose up --build
```

#### **Verify the server**

```bash
curl http://localhost:8080/api/v1/health
```

### Local Development Setup

#### **Install dependencies**

```bash
go mod download
```

#### **Start PostgreSQL**

```bash
# Using Docker
docker run --name pingpong_db \
  -e POSTGRES_USER=pingpong_user \
  -e POSTGRES_PASSWORD=pingpong_password_2024 \
  -e POSTGRES_DB=pingpong \
  -p 5432:5432 \
  -d postgres:15-alpine
```

#### **Initialize database**

```bash
psql -h localhost -U pingpong_user -d pingpong -f init.sql
```

#### **Configure .env**

```bash
cp .env.example .env
```

#### **Run the server**

```bash
go run main.go database.go models.go routes.go handlers.go
```

---

## 📡 API Endpoints

### Base URL

```bash
http://localhost:8080/api/v1
```

### Health Check

```bash
GET /health
```

Returns server status.

---

## 👥 Player Management

### Get all players
```
GET /players
```
**Response:**
```json
{
  "message": "Players fetched successfully",
  "data": [
    {
      "id": "uuid",
      "full_name": "Nguyễn Văn A",
      "birth_year": 1990,
      "phone": "0912345678",
      "cccd": "001234567890",
      "avatar_url": "https://...",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Get player by ID
```
GET /players/:id
```

### Create player
```
POST /players
Content-Type: application/json

{
  "full_name": "Nguyễn Văn A",
  "birth_year": 1990,
  "phone": "0912345678",
  "cccd": "001234567890",
  "avatar_url": "https://..."
}
```

### Update player
```
PUT /players/:id
Content-Type: application/json

{
  "full_name": "Updated Name",
  "birth_year": 1991,
  "phone": "0987654321",
  "is_active": true
}
```

### Upload avatar
```
POST /players/avatar/upload
Content-Type: multipart/form-data

- avatar: [file]
- user_id: [uuid]
```

---

## 🏆 Team Management

### Get all teams
```
GET /teams
```

### Create team
```
POST /teams
Content-Type: application/json

{
  "name": "Club A",
  "short_name": "CA",
  "logo_url": "https://..."
}
```

---

## 🎖️ Rank Management

### Get all ranks
```
GET /ranks
```
Returns all player rank levels (C3, C2, C1, B3, B2, B1, A3, A2, A1).

---

## 🏅 Season Management

### Get all seasons
```
GET /seasons
```

### Get season by ID
```
GET /seasons/:id
```

### Create season
```
POST /seasons
Content-Type: application/json

{
  "name": "Giải đấu tháng 1/2024",
  "year": 2024
}
```

### Get players in season
```
GET /seasons/:id/players
```
Returns all players registered for a season with their details.

### Add player to season
```
POST /seasons/:id/players
Content-Type: application/json

{
  "player_id": "uuid",
  "team_id": "uuid",
  "rank_id": "C2"
}
```

---

## 📊 Leaderboard

### Get season leaderboard
```
GET /seasons/:id/leaderboard
```
Returns ranked list of players for a season.

**Response:**
```json
{
  "message": "Leaderboard fetched successfully",
  "data": [
    {
      "rank": 1,
      "player_id": "uuid",
      "player_name": "Nguyễn Văn A",
      "rank_id": "B1",
      "rank_name": "Hạng B1",
      "team_id": "uuid",
      "team_name": "Club A",
      "accumulated_points": 1850.5,
      "status": "ACTIVE",
      "display_order": 1
    }
  ]
}
```

---

## 💰 Point Management

### Adjust player points
```
POST /points/adjust
Content-Type: application/json

{
  "player_season_id": "uuid",
  "delta_points": 100.5,
  "reason": "Thắng trận đấu vòng 1"
}
```

### Get point logs
```
GET /points/logs/:playerSeasonId
```
Returns audit trail of all point changes for a player.

---

## 🎯 Fixture Management

### Get fixtures for season
```
GET /seasons/:id/fixtures
```

### Create fixture
```
POST /fixtures
Content-Type: application/json

{
  "season_id": "uuid",
  "round": 1,
  "home_team_id": "uuid",
  "guest_team_id": "uuid"
}
```

---

## 🎮 Match Management

### Get matches for fixture
```
GET /fixtures/:id/matches
```

### Create match
```
POST /matches
Content-Type: application/json

{
  "fixture_id": "uuid",
  "match_order": 1,
  "match_type": "SINGLE",
  "home_player1_id": "uuid",
  "home_player2_id": null,
  "guest_player1_id": "uuid",
  "guest_player2_id": null,
  "handicap_snapshot": null
}
```

### Record match result
```
POST /matches/:id/result
Content-Type: application/json

{
  "home_sets": [11, 9],
  "guest_sets": [8, 6],
  "winner_team_id": "uuid"
}
```

---

## 🗄️ Database Schema

### Core Tables

#### ranks
```sql
- id (VARCHAR 10, PK)
- sort_order (INT)
- min_score (INT, nullable)
- max_score (INT, nullable)
- standard_score (INT)
- description (TEXT)
```

#### teams
```sql
- id (UUID, PK)
- name (TEXT, UNIQUE)
- short_name (TEXT)
- logo_url (TEXT)
- created_at (TIMESTAMP)
```

#### players
```sql
- id (UUID, PK)
- full_name (TEXT)
- birth_year (INT)
- phone (TEXT)
- cccd (TEXT)
- avatar_url (TEXT)
- is_active (BOOLEAN)
- created_at (TIMESTAMP)
```

#### seasons
```sql
- id (UUID, PK)
- name (TEXT)
- year (INT)
- is_active (BOOLEAN)
- created_at (TIMESTAMP)
```

### Season Tables

#### player_seasons (⭐ Core)
```sql
- id (UUID, PK)
- season_id (UUID, FK)
- player_id (UUID, FK)
- team_id (UUID, FK)
- rank_id (VARCHAR 10, FK)
- accumulated_points (NUMERIC)
- status (TEXT)
- display_order (INT)
- created_at (TIMESTAMP)
- UNIQUE (season_id, player_id)
```

#### player_point_logs
```sql
- id (UUID, PK)
- player_season_id (UUID, FK)
- delta_points (NUMERIC)
- reason (TEXT)
- source (TEXT) -- MATCH, ADMIN_ADJUST, PENALTY, BONUS
- ref_id (UUID)
- created_at (TIMESTAMP)
```

### Competition Tables

#### fixtures
```sql
- id (UUID, PK)
- season_id (UUID, FK)
- round (INT)
- home_team_id (UUID, FK)
- guest_team_id (UUID, FK)
- home_score (INT)
- guest_score (INT)
- status (TEXT) -- SCHEDULED, ONGOING, COMPLETED
- created_at (TIMESTAMP)
```

#### matches
```sql
- id (UUID, PK)
- fixture_id (UUID, FK)
- match_order (INT)
- match_type (TEXT) -- SINGLE, DOUBLE
- home_player1_id (UUID, FK)
- home_player2_id (UUID, FK, nullable)
- guest_player1_id (UUID, FK)
- guest_player2_id (UUID, FK, nullable)
- handicap_snapshot (TEXT)
- home_sets (INT[])
- guest_sets (INT[])
- winner_team_id (UUID, FK, nullable)
- created_at (TIMESTAMP)
```

---

## 🔐 Data Principles

✅ **DO:**
- Always log point changes in `player_point_logs`
- Use transactions for related operations
- Store data at player_seasons level during active season
- Keep audit trail forever

❌ **DON'T:**
- Update `players.accumulated_points` directly (doesn't exist)
- Delete historical match data
- Modify player_seasons without logging changes
- Update points without audit trail

---

## 🛠️ Development

### Building from source
```bash
go build -o pingpong-backend .
./pingpong-backend
```

### Running tests
```bash
go test ./...
```

### Linting
```bash
go fmt ./...
golint ./...
```

---

## 🐳 Docker Commands

```bash
# Build image
docker build -t pingpong-backend:latest .

# Start services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down

# Clean up everything
docker-compose down -v
```

---

## 📝 Environment Variables

See `.env.example` for complete list.

Key variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `PORT` (default: 8080)
- `GIN_MODE` (debug/release)

---

## 🚨 Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Invalid request body | Malformed JSON or missing required fields |
| 404 | Not found | Resource doesn't exist |
| 500 | Internal server error | Database or server error |

---

## 📞 Support

For issues or questions:
1. Check database logs: `docker-compose logs postgres`
2. Check backend logs: `docker-compose logs backend`
3. Review database schema: `init.sql`
4. Check API documentation above

---

## 📄 License

This project is part of Hanoi Super League Ping Pong Tournament system.
