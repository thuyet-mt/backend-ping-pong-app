Dưới đây là **tài liệu thiết kế DB hoàn chỉnh (Final Version)** cho hệ thống League Pingpong, đã bao gồm:

* Multi-season
* Transfer giữa đội
* Snapshot ranking theo round
* Rating (ELO)
* Buffer ±25 chống nhảy hạng
* Snapshot standings tối ưu performance

Tài liệu này có thể cho team BE Go để thiết kế API.

---

# I. Tổng Quan Kiến Trúc

Hệ thống được thiết kế theo mô hình:

> Season-based League System với Snapshot Ranking

Mỗi mùa (Season) là một không gian dữ liệu độc lập.

Hệ thống chia thành 6 domain:

1. Season Management
2. Team Management
3. Player Management
4. Match Engine
5. Ranking & Rating Engine
6. Configuration (Rank System)

---

# II. ERD Logic (Quan hệ tổng thể)

```markdown
seasons
 ├── teams
 │     └── team_members
 │
 ├── rounds
 │
 ├── matches
 │     └── sub_matches
 │           └── match_player_points
 │
 ├── player_round_points
 │     └── player_round_standings
 │
 ├── player_ratings
 │
 └── team_standings

players
 ├── player_ratings
 ├── player_round_points
 ├── player_round_standings
 └── player_rank_history

rank_configs
```

---

# III. Domain Chi Tiết

---

# 1️⃣ SEASON DOMAIN

## seasons

Mỗi record đại diện một mùa giải.

| Cột        | Mô tả                        |
| ---------- | ---------------------------- |
| id         | UUID                         |
| name       | Tên mùa                      |
| start_date | Ngày bắt đầu                 |
| end_date   | Ngày kết thúc                |
| status     | UPCOMING / ACTIVE / FINISHED |

Rule:

* Chỉ có 1 season ACTIVE tại một thời điểm (enforce ở BE).

---

## rounds

Quản lý trạng thái từng vòng.

| Cột          | Mô tả                     |
| ------------ | ------------------------- |
| season_id    | FK                        |
| round_number | Số vòng                   |
| status       | OPEN / LOCKED / FINALIZED |

Rule:

* Không cho sửa match nếu round = FINALIZED
* FINALIZED chỉ chạy 1 lần
* UNIQUE(season_id, round_number)

---

# 2️⃣ TEAM DOMAIN

## teams

Đội thi đấu trong một season.

| Cột       | Mô tả   |
| --------- | ------- |
| id        | UUID    |
| season_id | FK      |
| name      | Tên đội |

UNIQUE(season_id, name)

---

## team_members

Quan hệ player – team theo mùa.

| Cột           | Mô tả              |
| ------------- | ------------------ |
| team_id       | FK                 |
| player_id     | FK                 |
| joined_round  | Round bắt đầu      |
| left_round    | Round rời          |
| transfer_type | INITIAL / TRANSFER |

Rule:

* joined_round > 0
* left_round >= joined_round

Cho phép:

* Transfer giữa vòng
* Player chơi nhiều season

---

# 3️⃣ PLAYER DOMAIN

## players

Thông tin cố định, không phụ thuộc season.

| Cột           | Mô tả     |
| ------------- | --------- |
| id            | UUID      |
| full_name     | Tên       |
| date_of_birth | Ngày sinh |
| gender        | Giới tính |

---

# 4️⃣ RATING & RANK SYSTEM

---

## rank_configs

Định nghĩa hệ thống hạng và buffer ±25.

| Cột               | Mô tả          |
| ----------------- | -------------- |
| rank              | F1, E2…        |
| min_points        | Điểm tối thiểu |
| max_points        | Điểm tối đa    |
| promotion_buffer  | Mặc định 25    |
| relegation_buffer | Mặc định 25    |

Rank không còn phụ thuộc trực tiếp vào điểm.
Nó là một state machine.

---

## player_ratings

Trạng thái player theo từng season.

| Cột                | Mô tả           |
| ------------------ | --------------- |
| player_id          | FK              |
| season_id          | FK              |
| points             | Rating hiện tại |
| rank               | Rank hiện tại   |
| accumulated_points | Tổng điểm mùa   |
| updated_at         | Timestamp       |

UNIQUE(player_id, season_id)

Đây là bảng để load leaderboard nhanh.

---

## player_rank_history

Lưu lịch sử lên / xuống hạng.

| Cột          | Mô tả         |
| ------------ | ------------- |
| player_id    | FK            |
| season_id    | FK            |
| round_number | Vòng thay đổi |
| old_rank     | Hạng cũ       |
| new_rank     | Hạng mới      |

---

# 5️⃣ MATCH ENGINE

---

## matches

Trận đấu giữa 2 đội.

| Cột          | Mô tả                           |
| ------------ | ------------------------------- |
| season_id    | FK                              |
| round_number | Round                           |
| home_team_id | FK                              |
| away_team_id | FK                              |
| home_points  | Điểm đội                        |
| away_points  | Điểm đội                        |
| status       | PENDING / COMPLETED / FINALIZED |

Rule:

* home_team_id != away_team_id

---

## sub_matches

Tối đa 9 trận con.

| Cột            | Mô tả           |
| -------------- | --------------- |
| match_order    | 1–9             |
| match_type     | SINGLE / DOUBLE |
| home_player1   | FK              |
| home_player2   | FK              |
| away_player1   | FK              |
| away_player2   | FK              |
| home_sets      |                 |
| away_sets      |                 |
| winner_team_id | FK              |

UNIQUE(match_id, match_order)

---

## match_player_points

Delta rating (ELO change).

| Cột          | Mô tả |
| ------------ | ----- |
| sub_match_id | FK    |
| player_id    | FK    |
| delta_points | +/-   |

UNIQUE(sub_match_id, player_id)

---

# 6️⃣ ROUND POINT SYSTEM (Season Points)

---

## player_round_points

Điểm tích lũy từng vòng.

| Cột           | Mô tả        |
| ------------- | ------------ |
| player_id     | FK           |
| season_id     | FK           |
| round_number  |              |
| points_earned | Điểm vòng đó |

UNIQUE(player_id, season_id, round_number)

---

## player_round_standings

Snapshot bảng xếp hạng sau từng vòng.

| Cột                | Mô tả             |
| ------------------ | ----------------- |
| accumulated_points | Tổng đến round đó |
| rank_position      | Thứ hạng          |

UNIQUE(player_id, season_id, round_number)

Dùng để:

* Xem lịch sử thứ hạng
* Vẽ biểu đồ progression

---

## team_standings

Bảng xếp hạng đội.

| Cột          | Mô tả    |
| ------------ | -------- |
| total_points | Điểm đội |
| match_wins   |          |
| set_wins     |          |

UNIQUE(season_id, team_id)

---

# IV. Luồng Nghiệp Vụ Chuẩn

---

# 1️⃣ Nhập kết quả trận

Insert:

* sub_matches
* match_player_points

Update:

* matches.status = COMPLETED

---

# 2️⃣ Finalize Round (Transaction bắt buộc)

Flow:

1. Kiểm tra round đang OPEN
2. Update round → FINALIZED
3. Insert player_round_points
4. Snapshot player_round_standings (RANK() window function)
5. Update player_ratings.accumulated_points
6. Update rating points
7. Tính rank mới theo buffer ±25
8. Insert player_rank_history nếu đổi hạng
9. Update team_standings

Toàn bộ phải chạy trong 1 transaction.

---

# V. Rank Buffer ±25 (State Machine)

Rule:

Promotion:
points >= max_points + promotion_buffer

Relegation:
points < min_points - relegation_buffer

Rank không được tính trực tiếp từ khoảng điểm.

---

# VI. Performance Strategy

* Không SUM realtime
* Snapshot sau mỗi round
* Index quan trọng:

```sql
INDEX(player_id, season_id)
INDEX(season_id, round_number)
INDEX(match_id)
INDEX(team_id)
```

---

# VII. API Design Mapping (Cho BE Go)

---

## Season

POST /seasons
POST /seasons/{id}/start
POST /seasons/{id}/finish

---

## Teams

POST /seasons/{id}/teams
POST /teams/{id}/add-player
POST /teams/{id}/transfer

---

## Matches

POST /matches
POST /matches/{id}/sub-match
POST /matches/{id}/complete

---

## Round Control

POST /seasons/{id}/rounds/{round}/finalize

---

## Leaderboard

GET /seasons/{id}/leaderboard
GET /seasons/{id}/rounds/{round}/leaderboard
GET /players/{id}/history

---

# VIII. Kiến Trúc Đạt Được

Hệ thống hiện tại:

* Multi-season
* Transfer giữa đội
* Rating ELO
* Buffer chống nhảy hạng
* Snapshot ranking từng vòng
* Snapshot ranking toàn mùa
* Audit lịch sử lên/xuống hạng
* Tối ưu performance cho production
