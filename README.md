# 🏓 Ping Pong Tournament Backend

Hệ thống quản lý giải đấu bóng bàn toàn diện, được xây dựng bằng **Go** và **PostgreSQL**. Cung cấp API RESTful hoàn chỉnh với kiến trúc lớp 3 (Repository, Service, Handler), xử lý điểm số, quản lý bảng xếp hạng, và tính năng nâng cao như snapshot trận đấu, materialized views cho hiệu suất cao.

**Phiên bản**: 2.0.0 | **Status**: Production-Ready ✅

## 📋 Mục lục

- Features (Tính năng chính)
- System Requirements (Yêu cầu hệ thống)
- Project Structure (Cấu trúc dự án)
- Installation (Cài đặt)
- Running (Chạy ứng dụng)
- API Documentation
- Architecture (Kiến trúc hệ thống)
- Error Handling (Xử lý lỗi)
- Troubleshooting

## ✨ Tính năng chính


### Quản lý VĐV & Đội

- ✅ Tạo, cập nhật, liệt kê VĐV
- ✅ Quản lý đội bóng
- ✅ Chỉ định VĐV vào mùa giải
- ✅ Điều chỉnh xếp hạng (Rank)

### Quản lý Giải đấu & Trận đấu

- ✅ Tạo mùa giải (Season)
- ✅ Quản lý bảng (Fixture) và trận đấu (Match)
- ✅ Ghi nhận kết quả trận đấu
- ✅ Snapshot trạng thái tại thời điểm trận đấu

### Bảng Xếp Hạng & Điểm

- ✅ Bảng xếp hạng theo mùa giải (Leaderboard)
- ✅ Quản lý điểm (Points) với audit log
- ✅ Materialized view cho hiệu suất cao (300-500 VĐV)
- ✅ Lịch sử thay đổi điểm

### Tính Năng Nâng Cao

- ✅ Xử lý lỗi chuẩn hóa với error codes cho ứng dụng di động
- ✅ Architecture 3-lớp: Handler → Service → Repository
- ✅ Background Jobs cho tác vụ tự động (Rank recalc, Leaderboard refresh)
- ✅ Hỗ trợ tải ảnh VĐV (Avatar)
- ✅ CORS middleware cho frontend integration

## 🖥️ Yêu cầu hệ thống

- **Go**: 1.21+
- **PostgreSQL**: 15+
- **Docker & Docker Compose**: (Tuỳ chọn)
- **RAM**: 2GB | **CPU**: 2 cores | **Disk**: 500MB

## 📁 Cấu trúc dự án

```bash
backend-ping-pong-app/
├── cmd/
│   └── server/
│       └── main.go                # Entry point duy nhất
│
├── internal/
│   ├── config/                    # Load env, config
│   │   └── config.go
│   │
│   ├── database/                  # DB connection, migration helper
│   │   ├── postgres.go
│   │   └── tx.go
│   │
│   ├── middleware/                # Gin middleware
│   │   ├── auth_firebase.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recovery.go
│   │
│   ├── errors/                    # App error chuẩn hoá
│   │   ├── errors.go
│   │   └── http_mapper.go
│   │
│   ├── models/                    # Domain models (DB + DTO)
│   │   ├── player.go
│   │   ├── team.go
│   │   ├── season.go
│   │   ├── match.go
│   │   ├── rank.go
│   │   └── dto.go
│   │
│   ├── repository/                # DB access layer
│   │   ├── player_repo.go
│   │   ├── team_repo.go
│   │   ├── season_repo.go
│   │   ├── match_repo.go
│   │   ├── rank_repo.go
│   │   └── repository.go          # interface aggregator
│   │
│   ├── service/                   # Business logic
│   │   ├── player_service.go
│   │   ├── team_service.go
│   │   ├── season_service.go
│   │   ├── match_service.go
│   │   ├── rank_service.go
│   │   └── service.go             # Service container
│   │
│   ├── handlers/                  # HTTP handlers
│   │   ├── player_handler.go
│   │   ├── team_handler.go
│   │   ├── season_handler.go
│   │   ├── match_handler.go
│   │   ├── rank_handler.go
│   │   ├── upload_handler.go
│   │   ├── routes.go
│   │   └── handlers.go            # Handler container
│   │
│   ├── jobs/                      # Background jobs
│   │   ├── leaderboard_job.go
│   │   ├── season_reset_job.go
│   │   └── jobs.go
│   │
│   ├── utils/                     # Helper chung
│   │   ├── time.go
│   │   ├── uuid.go
│   │   └── pagination.go
│   │
│   └── constants/                 # Enum, constant
│       └── rank.go
│
├── sql/
│   ├── init.sql                   # Init schema
│   └── seed.sql                   # Seed data (optional)
│
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── docs/
│   ├── api.md
│   ├── database.md
│   └── deployment.md
│
├── .env.example
├── Makefile
├── go.mod
├── go.sum
└── README.md

```

## 🚀 Cài đặt & Chạy


### Quick Start (Docker)

```bash
# Clone repo
git clone <repository-url>
cd backend-ping-pong-app

# Copy environment
cp config/.env.example .env

# Start services
docker-compose -f docker/docker-compose.yml up --build

# Test
curl http://localhost:8080/api/v1/health
```

### Local Development

```bash
# Setup
go mod download
go mod tidy

# Configure
cp config/.env.example .env
# Edit .env with your database info

# Create PostgreSQL database
createdb pingpong
psql -d pingpong -f sql/init.sql

# Build & Run
go build -o app cmd/server/main.go
./app

# Or use make
make run
```

## 📡 API Base URL

```
http://localhost:8080/api/v1
```

## 🏗️ Kiến trúc 3-Lớp

```plaintext
Handlers (HTTP)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
PostgreSQL Database
```

### Layer Details


**Handlers** (`internal/handlers/`)

- Xử lý HTTP requests
- Validate input
- Call service methods
- Return JSON responses

**Service** (`internal/service/`)

- Business logic & rules
- Validation
- Point calculations
- Error handling

**Repository** (`internal/repository/`)

- Pure database queries
- No business logic
- Connection pooling
- Transaction support

## 🗄️ Database Schema

**Key Tables:**

- `players` - VĐV
- `teams` - Đội
- `seasons` - Mùa giải
- `player_seasons` - Tham gia VĐV-Mùa
- `ranks` - Xếp hạng
- `fixtures` - Bảng
- `matches` - Trận đấu
- `player_point_logs` - Lịch sử điểm (Audit)

**Views:**

- `v_season_leaderboard` - Real-time bảng xếp
- `mat_season_leaderboard` - Materialized view (refreshed 5min)

## 📊 Main Endpoints

```bash
GET    /players                    # Danh sách VĐV
POST   /players                    # Tạo VĐV
GET    /seasons                    # Danh sách mùa giải
POST   /seasons                    # Tạo mùa giải
GET    /seasons/:id/players        # VĐV trong mùa
GET    /seasons/:id/leaderboard    # Bảng xếp hạng
POST   /points/adjust              # Điều chỉnh điểm
GET    /points/logs/:id            # Lịch sử điểm
POST   /matches/:id/result         # Ghi nhận kết quả
```

Xem [docs/API.md](docs/API.md) để chi tiết.

## ❌ Error Handling

All errors return standardized structure:

```json
{

  "code": "ERROR_CODE",
  "message": "Mô tả lỗi",
  "status_code": 400,
  "details": {}
}
```

Common error codes:
- `ErrorPlayerNotFound` (404)
- `ErrorPlayerAlreadyInSeason` (409)
- `ErrorMatchAlreadyRecorded` (409)
- `ErrorNegativePointsResult` (400)
- `InvalidInput` (400)

## 🔧 Commands

```bash
# Build
make build

# Run
make run

# Docker
docker-compose -f docker/docker-compose.yml up -d

# Test
curl http://localhost:8080/api/v1/health
```

## 🧪 Testing với Postman

```bash
Base URL: http://localhost:8080/api/v1

# 1. Health Check
GET /health

# 2. Tạo VĐV
POST /players
{
  "full_name": "Người chơi A",
  "birth_year": 1995
}

# 3. Danh sách VĐV
GET /players
```

## 🐳 Docker Deployment

```bash
# Build
docker build -t pingpong-backend docker/

# Run
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=pingpong \
  -e DB_PASSWORD=password \
  pingpong-backend

# Compose
docker-compose -f docker/docker-compose.yml up -d
```

## 📚 Background Jobs

1. **RankRecalcJob** (mỗi giờ) - Tính toán lại xếp hạng


2. **LeaderboardMatJob** (mỗi 5 phút) - Refresh materialized view


3. **AnomalyDetectionJob** (mỗi 6 giờ) - Phát hiện hành vi bất thường


4. **StatsAggregationJob** - Thu thập thống kê


## 🆘 Troubleshooting

### Database Connection Error
```bash
# Check PostgreSQL
psql -U postgres

# Check config
cat .env | grep DB_

# Reset database
psql -d pingpong -f sql/init.sql
```

### Port Already in Use
```bash
# Change in .env
PORT=8081

# Or kill process
lsof -i :8080 | grep -v COMMAND | awk '{print $2}' | xargs kill
```

### Docker Issues
```bash
# View logs
docker-compose -f docker/docker-compose.yml logs -f

# Rebuild
docker-compose -f docker/docker-compose.yml up --build

# Clean
docker-compose -f docker/docker-compose.yml down -v
```

## 📖 Thêm tài liệu

- [API Reference](docs/API.md)
- [Database Schema](docs/DATABASE.md)
- [Deployment Guide](docs/DEPLOYMENT.md)

## 🤝 Đóng góp

Fork → Feature Branch → Commit → Push → Pull Request

## 📄 License

MIT License

---

## Credits

Built with ❤️ using Go + PostgreSQL

For issues: GitHub Issues | Support: [email/contact]
