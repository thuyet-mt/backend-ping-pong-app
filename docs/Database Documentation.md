# 📘 TÀI LIỆU CƠ SỞ DỮ LIỆU – HANOI SUPER LEAGUE (PINGPONG)

Tài liệu này mô tả **đầy đủ kiến trúc database**, cách **khởi tạo**, **cấu hình**, và **viết query chuẩn** cho hệ thống giải đấu bóng bàn, phục vụ **Flutter App (FE)** và **Go Backend (BE)**.

---

## 1. Tổng quan kiến trúc

### 🎯 Mục tiêu thiết kế

* Scale tốt (nhiều mùa giải, nhiều VĐV)
* Không mất lịch sử (audit & analytics)
* Phù hợp giải truyền thống + BTC can thiệp
* Query đơn giản, rõ ràng

### 🧱 Phân tầng dữ liệu

| Tầng        | Mô tả                   |
| ----------- | ----------------------- |
| Core        | players, teams, ranks   |
| Season      | seasons, player_seasons |
| Competition | fixtures, matches       |
| Analytics   | player_point_logs       |
| Import      | staging_players         |

---

## 2. Khởi tạo Database

### 2.1 Tạo database

```sql
CREATE DATABASE pingpong
  WITH OWNER = postgres
  ENCODING = 'UTF8'
  LC_COLLATE = 'en_US.UTF-8'
  LC_CTYPE = 'en_US.UTF-8';
```

### 2.2 Extension bắt buộc

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

---

## 3. Các bảng chính

### 3.1 ranks – Hạng trình độ

```sql
CREATE TABLE ranks (
  id VARCHAR(10) PRIMARY KEY,
  sort_order INT NOT NULL,
  min_score INT,
  max_score INT,
  standard_score INT,
  description TEXT
);
```

Dùng cho:

* Tính chấp
* Xếp hạng
* Analytics

---

### 3.2 teams – Đội bóng

```sql
CREATE TABLE teams (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  short_name TEXT,
  logo_url TEXT,
  created_at TIMESTAMP DEFAULT now()
);
```

---

### 3.3 players – Vận động viên (master)

```sql
CREATE TABLE players (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name TEXT NOT NULL,
  birth_year INT,
  phone TEXT,
  cccd TEXT,
  avatar_url TEXT,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT now()
);
```

⚠️ **Không lưu điểm, không lưu đội, không lưu mùa giải**

---

### 3.4 seasons – Mùa giải

```sql
CREATE TABLE seasons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  year INT,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT now()
);
```

---

### 3.5 player_seasons – VĐV theo mùa giải (QUAN TRỌNG)

```sql
CREATE TABLE player_seasons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  season_id UUID REFERENCES seasons(id),
  player_id UUID REFERENCES players(id),
  team_id UUID REFERENCES teams(id),
  rank_id VARCHAR(10) REFERENCES ranks(id),
  accumulated_points NUMERIC(10,2) DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  display_order INT,
  UNIQUE (season_id, player_id)
);
```

👉 **Tất cả điểm, hạng, trạng thái đều nằm ở đây**

---

### 3.6 player_point_logs – Nhật ký điểm (AUDIT)

```sql
CREATE TABLE player_point_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_season_id UUID REFERENCES player_seasons(id),
  delta_points NUMERIC(10,2),
  reason TEXT,
  source TEXT, -- MATCH, ADMIN_ADJUST
  ref_id UUID,
  created_at TIMESTAMP DEFAULT now()
);
```

Dùng cho:

* Truy vết gian lận
* Rollback
* Thống kê

---

### 3.7 fixtures – Đối đầu CLB

```sql
CREATE TABLE fixtures (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  season_id UUID REFERENCES seasons(id),
  round INT,
  home_team_id UUID REFERENCES teams(id),
  guest_team_id UUID REFERENCES teams(id),
  home_score INT DEFAULT 0,
  guest_score INT DEFAULT 0,
  status TEXT DEFAULT 'SCHEDULED'
);
```

---

### 3.8 matches – Trận đấu con

```sql
CREATE TABLE matches (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  fixture_id UUID REFERENCES fixtures(id) ON DELETE CASCADE,
  match_order INT,
  match_type TEXT, -- SINGLE / DOUBLE
  home_player1_id UUID REFERENCES players(id),
  home_player2_id UUID REFERENCES players(id),
  guest_player1_id UUID REFERENCES players(id),
  guest_player2_id UUID REFERENCES players(id),
  handicap_snapshot TEXT,
  home_sets INT[],
  guest_sets INT[],
  winner_team_id UUID REFERENCES teams(id)
);
```

---

## 4. Import dữ liệu

### 4.1 staging_players (import tạm)

```sql
CREATE TABLE staging_players (
  vdv_ten TEXT,
  nam_sinh INT,
  vdv_hang TEXT,
  diem_tich_luy NUMERIC,
  doi_bong_ten TEXT,
  mua_giai_ten TEXT,
  trang_thai_thi_dau TEXT,
  stt INT
);
```

### 4.2 Import vào player_seasons

```sql
INSERT INTO player_seasons (...)
SELECT ...
FROM staging_players s
JOIN players p ...
JOIN seasons se ...
LEFT JOIN teams t ...
ON CONFLICT (season_id, player_id) DO NOTHING;
```

---

## 5. Query mẫu thường dùng

### BXH mùa giải

```sql
SELECT p.full_name, ps.accumulated_points, r.id AS rank
FROM player_seasons ps
JOIN players p ON p.id = ps.player_id
JOIN ranks r ON r.id = ps.rank_id
WHERE ps.season_id = :season_id
ORDER BY ps.accumulated_points DESC;
```

### Cộng điểm + log

```sql
UPDATE player_seasons
SET accumulated_points = accumulated_points + :delta
WHERE id = :ps_id;

INSERT INTO player_point_logs (...);
```

---

## 6. Nguyên tắc vận hành

* ❌ Không update trực tiếp players.accumulated_points
* ✅ Mọi thay đổi điểm phải có log
* ✅ Rank có thể override bởi BTC
* ✅ Không xóa dữ liệu đã thi đấu

---

## 7. Phù hợp cho Flutter + Go

* Flutter: chỉ gọi API, không logic điểm
* Go: service layer tính điểm + transaction
* DB: source of truth

---

## 8. Hướng mở rộng

* ELO rating
* Playoff
* Multiple division
* Sponsor / MVP analytics

---

📌 **Tài liệu này đủ dùng cho production & mở rộng lâu dài**
