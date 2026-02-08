# 🚀 Deployment & Architecture Guide

## System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Flutter Mobile App                   │
│                     (Frontend Client)                    │
└────────────────────────┬────────────────────────────────┘
                         │
                    HTTP/HTTPS
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│              Go Backend API           │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Main Server:                                    │  │
│  │  - REST API Endpoints (/api/v1/...)            │  │
│  │  - Request Validation                           │  │
│  │  - Error Handling                               │  │
│  │  - CORS Support                                 │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Database Layer:                                 │  │
│  │  - Query Building                               │  │
│  │  - Transaction Management                       │  │
│  │  - Connection Pooling                           │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Business Logic:                                 │  │
│  │  - Point Calculation                            │  │
│  │  - Ranking Algorithm                            │  │
│  │  - Audit Logging                                │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
                   TCP Port 5432
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                PostgreSQL Database                       │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Core Data:                                      │  │
│  │  - Players, Teams, Seasons, Ranks               │  │
│  │  - Player_Seasons (many-to-many)                │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Competition Data:                               │  │
│  │  - Fixtures (Team matches)                       │  │
│  │  - Matches (Individual games)                    │  │
│  │  - Match Results & Scoring                       │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Audit & Analytics:                              │  │
│  │  - Point Logs (100% audit trail)                │  │
│  │  - Leaderboards & Rankings                       │  │
│  │  - Statistical Views                             │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## 📈 Data Flow Example: Recording Match Result

```
1. Flutter App sends match result
   POST /api/v1/matches/{id}/result
   {
     "home_sets": [11, 9],
     "guest_sets": [8, 6],
     "winner_team_id": "uuid"
   }

2. API Handler receives & validates request
   - Check match exists
   - Validate data
   - Parse request body

3. Database Layer executes transaction
   Step 1: Update match record
           SET home_sets = [11, 9]
           SET guest_sets = [8, 6]
           SET winner_team_id = uuid

   Step 2: Identify winning players from fixture
           Query match to get player IDs

   Step 3: Calculate points for each player
           - Winning team: +100 points
           - Losing team: +50 points
           (Can be customized)

   Step 4: Update player_seasons
           UPDATE player_seasons
           SET accumulated_points = accumulated_points + delta
           WHERE player in (winner1, winner2)

   Step 5: Create audit logs
           INSERT INTO player_point_logs
           FOR EACH player
           WITH source = "MATCH"
           WITH ref_id = match_id

4. Response sent to client
   {
     "message": "Match result recorded successfully",
     "data": { match object }
   }

5. Flutter App updates UI
   - Refresh leaderboard
   - Show updated points
   - Display match result
```

---

## 🐳 Docker Deployment Layers

### Layer 1: PostgreSQL Container
```dockerfile
FROM postgres:15-alpine

VOLUMES:
  - postgres_data (persistent)
  
ENTRYPOINT:
  - init.sql (auto-run on first start)
  
NETWORK:
  - pingpong_network (internal)

EXPOSED:
  - Port 5432 (internal only)
```

### Layer 2: Go Backend Container
```dockerfile
FROM golang:1.21-alpine (build stage)
  - Compile Go application
  - Optimize for production
  
FROM alpine:latest (runtime stage)
  - Minimal image size
  - Only runtime binary
  
ENTRYPOINT:
  - ./pingpong-backend
  
NETWORKS:
  - Connects to postgres container
  
EXPOSED:
  - Port 8080 (to host)
```

### Docker Compose Orchestration
```yaml
Network Setup:
  - Creates pingpong_network bridge
  - Allows container-to-container DNS resolution
  - Isolates from other containers

Volume Management:
  - postgres_data: Persistent database storage
  - uploads/: Shared file storage

Health Checks:
  - PostgreSQL: pg_isready health check
  - Backend: HTTP /health endpoint

Init Sequence:
  1. Start PostgreSQL
  2. Wait for postgres health check (ready)
  3. Start Backend (depends_on: postgres)
  4. Backend connects to postgres:5432
  5. init.sql runs once on first start
```

---

## 🔧 Configuration Management

### Environment Variables Priority

```
1. .env file (local development)
   ↓ overrides
2. Docker compose environment section
   ↓ overrides
3. System environment variables
   ↓ overrides
4. Hard-coded defaults in code
```

### Key Configuration Points

```go
// Database Connection
{
  host: env("DB_HOST", "localhost")
  port: env("DB_PORT", "5432")
  user: env("DB_USER", "pingpong_user")
  password: env("DB_PASSWORD", "***")
  database: env("DB_NAME", "pingpong")
  sslmode: "disable" (local), "require" (production)
}

// Server
{
  port: env("PORT", "8080")
  gin_mode: env("GIN_MODE", "debug")
}

// Files
{
  upload_dir: env("UPLOAD_DIR", "./uploads")
  max_size: env("MAX_UPLOAD_SIZE", "5242880")
}
```

---

## 📋 Deployment Checklist

### Pre-Deployment
- [ ] Review all database migrations
- [ ] Backup existing database (if upgrading)
- [ ] Test locally with docker-compose
- [ ] Verify all environment variables
- [ ] Check disk space requirements
- [ ] Review security settings
- [ ] Test API endpoints

### Deployment
- [ ] Build Docker image: `docker build -t app:version .`
- [ ] Tag image: `docker tag app:version registry/app:version`
- [ ] Push to registry: `docker push registry/app:version`
- [ ] Update docker-compose.yml with new version
- [ ] Backup current database
- [ ] Deploy: `docker-compose up -d`
- [ ] Verify services started: `docker ps`
- [ ] Health check: `curl http://localhost:8080/api/v1/health`
- [ ] Run smoke tests
- [ ] Monitor logs: `docker logs -f backend`

### Post-Deployment
- [ ] Monitor application logs
- [ ] Check database connections
- [ ] Verify API response times
- [ ] Test critical workflows
- [ ] Confirm backup created
- [ ] Document any issues
- [ ] Plan rollback if needed

---

## 🔐 Security Considerations

### Database Security
```sql
-- Connection security
- Use SSL/TLS in production
- Change default passwords immediately
- Use strong passwords (20+ chars, mixed)
- Limit database user permissions

-- Data security
- Enable row-level security if needed
- Regular backups (daily minimum)
- Backup encryption
- Off-site backup storage

-- Audit trail
- All modifications logged
- Cannot delete audit logs
- Timestamp all changes
- Track user/admin actions
```

### API Security
```go
-- Input validation
- Validate all request data
- Check UUID formats
- Validate array lengths
- Sanitize string inputs

-- Error handling
- Don't expose internal errors to client
- Log detailed errors server-side
- Generic error messages to user
- Monitor error patterns

-- Rate limiting
- Limit API calls per IP (future)
- Implement CORS properly
- Validate JWT tokens (future)
- Monitor for suspicious activity
```

### File Upload Security
```
- Validate MIME types
- Check file extensions
- Scan for malware (future)
- Limit file size
- Use randomized filenames
- Store outside public folder
```

---

## 📊 Performance Optimization

### Database Indexing
```sql
-- Existing indexes:
CREATE INDEX idx_player_seasons_season_id
CREATE INDEX idx_player_seasons_player_id
CREATE INDEX idx_player_seasons_accumulated_points
CREATE INDEX idx_point_logs_player_season_id
CREATE INDEX idx_point_logs_created_at
CREATE INDEX idx_fixtures_season_id
CREATE INDEX idx_fixtures_round
CREATE INDEX idx_matches_fixture_id

-- Consider adding for future:
- player_seasons(team_id)
- player_seasons(status)
- player_point_logs(source)
- fixtures(status)
```

### Query Optimization
```go
// Good
- Use joins instead of N+1 queries
- Filter early in queries
- Select only needed columns
- Use indexed columns in WHERE
- Cache leaderboards (materialized view)

// Bad
- Load all players then filter in Go
- Nested loops in application code
- SELECT * (get unnecessary columns)
- No indexes on filter columns
```

### Caching Strategy (Future)
```
- Leaderboards: Cache 5 minutes
- Player profiles: Cache 1 hour
- Team data: Cache 1 day
- Invalidate on data change
- Use Redis for distributed cache
```

---

## 🔄 Backup & Recovery

### Automated Backups
```bash
#!/bin/bash
# daily-backup.sh
docker-compose exec postgres pg_dump \
  -U pingpong_user pingpong > backup_$(date +%Y%m%d).sql

# Or using Docker volume mount
tar czf backup_$(date +%Y%m%d).tar.gz /var/lib/postgresql/data
```

### Recovery Procedure
```bash
# 1. Stop services
docker-compose down

# 2. Restore database
docker-compose up -d postgres && sleep 5
docker-compose exec -T postgres psql \
  -U pingpong_user pingpong < backup_20240101.sql

# 3. Verify restore
docker-compose exec postgres psql \
  -U pingpong_user pingpong -c "SELECT COUNT(*) FROM players;"

# 4. Start backend
docker-compose up -d backend
```

### Point-in-Time Recovery
```
Using PostgreSQL WAL (Write-Ahead Log):
- Enable archiving in postgresql.conf
- Archive WAL files to backup storage
- Can restore to specific point in time
- Useful for large databases
```

---

## 📈 Monitoring & Logging

### Application Logs
```bash
# View logs
docker-compose logs -f backend

# View with timestamps
docker-compose logs --timestamps -f backend

# View last 100 lines
docker-compose logs --tail=100 backend
```

### Database Logs
```bash
# PostgreSQL logs
docker-compose logs -f postgres

# This is where connection issues appear
```

### Key Metrics to Monitor
```
- API Response Times
- Database Query Times (> 1s = slow)
- Point calculation accuracy (spot check)
- File uploads (success/failure rate)
- Leaderboard consistency
```

---

## 🚀 Production Deployment

### Recommended Setup

```
┌──────────────────────────────────────────┐
│    Cloudflare / Reverse Proxy            │
│    (SSL, Rate Limiting, DDoS)            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│  Load Balancer (if multi-instance)      │
│  (Nginx, HAProxy)                       │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
┌───────▼───────┐   ┌──────▼──────┐
│  Backend 1    │   │  Backend 2   │
│  Container    │   │  Container   │
└───────┬───────┘   └──────┬──────┘
        │                  │
        └────────┬─────────┘
                 │
        ┌────────▼────────┐
        │  PostgreSQL DB  │
        │  (Managed)      │
        └─────────────────┘
```

### Production Checklist
- [ ] Use managed PostgreSQL (AWS RDS, Azure Database)
- [ ] Enable automated backups
- [ ] Use SSL/TLS certificates
- [ ] Enable logging and monitoring
- [ ] Set up alerting
- [ ] Use environment-specific configs
- [ ] Load balance if multiple instances
- [ ] Use CDN for static files
- [ ] Regular security audits
- [ ] Performance testing

---

## 🆘 Troubleshooting

### Common Issues

#### Database Connection Refused
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker-compose logs postgres

# Verify network
docker network ls

# Restart
docker-compose restart postgres
```

#### Slow Queries
```bash
# Enable query logging
ALTER SYSTEM SET log_min_duration_statement = 1000; -- Log queries > 1s
SELECT pg_reload_conf();

# Check slow query log
docker-compose exec postgres \
  tail -f /var/log/postgresql/postgresql.log
```

#### Out of Disk Space
```bash
# Check usage
docker system df

# Clean unused volumes
docker volume prune

# Clean build cache
docker builder prune
```

---

**Production deployment requires careful planning and testing.**
**Always test migrations on a staging environment first.**
