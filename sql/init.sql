-- ==================== Extensions ====================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==================== Ranks Table ====================
CREATE TABLE IF NOT EXISTS ranks (
  id VARCHAR(10) PRIMARY KEY,
  sort_order INT NOT NULL UNIQUE,
  min_score INT,
  max_score INT,
  standard_score INT NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT now()
);

-- Seed initial ranks data
INSERT INTO ranks (id, sort_order, min_score, max_score, standard_score, description)
VALUES 
  ('C3', 1, NULL, 499, 50, 'Hạng C3 - Người mới'),
  ('C2', 2, 500, 799, 60, 'Hạng C2'),
  ('C1', 3, 800, 1099, 70, 'Hạng C1'),
  ('B3', 4, 1100, 1399, 80, 'Hạng B3'),
  ('B2', 5, 1400, 1699, 90, 'Hạng B2'),
  ('B1', 6, 1700, 1999, 100, 'Hạng B1'),
  ('A3', 7, 2000, 2499, 110, 'Hạng A3'),
  ('A2', 8, 2500, 2999, 120, 'Hạng A2'),
  ('A1', 9, 3000, NULL, 130, 'Hạng A1 - Xuất sắc')
ON CONFLICT (id) DO NOTHING;

-- ==================== Teams Table ====================
CREATE TABLE IF NOT EXISTS teams (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  short_name TEXT,
  logo_url TEXT,
  created_at TIMESTAMP DEFAULT now()
);

-- ==================== Players Table ====================
CREATE TABLE IF NOT EXISTS players (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name TEXT NOT NULL,
  birth_year INT,
  phone TEXT,
  cccd TEXT,
  avatar_url TEXT,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT now()
);

-- ==================== Seasons Table ====================
CREATE TABLE IF NOT EXISTS seasons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  year INT,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT now()
);

-- ==================== Player Seasons Table ====================
-- This is the core table storing player data during a season
CREATE TABLE IF NOT EXISTS player_seasons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
  player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  team_id UUID NOT NULL REFERENCES teams(id),
  rank_id VARCHAR(10) NOT NULL REFERENCES ranks(id),
  accumulated_points NUMERIC(10,2) DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  display_order INT,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  UNIQUE (season_id, player_id)
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_player_seasons_season_id ON player_seasons(season_id);
CREATE INDEX IF NOT EXISTS idx_player_seasons_player_id ON player_seasons(player_id);
CREATE INDEX IF NOT EXISTS idx_player_seasons_accumulated_points ON player_seasons(accumulated_points DESC);

-- ==================== Player Point Logs Table ====================
-- Audit trail for all point changes
CREATE TABLE IF NOT EXISTS player_point_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_season_id UUID NOT NULL REFERENCES player_seasons(id) ON DELETE CASCADE,
  delta_points NUMERIC(10,2) NOT NULL,
  reason TEXT,
  source TEXT DEFAULT 'MATCH', -- MATCH, ADMIN_ADJUST, PENALTY, BONUS
  ref_id UUID,
  created_at TIMESTAMP DEFAULT now()
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_point_logs_player_season_id ON player_point_logs(player_season_id);
CREATE INDEX IF NOT EXISTS idx_point_logs_created_at ON player_point_logs(created_at DESC);

-- ==================== Fixtures Table ====================
CREATE TABLE IF NOT EXISTS fixtures (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
  round INT NOT NULL,
  home_team_id UUID NOT NULL REFERENCES teams(id),
  guest_team_id UUID NOT NULL REFERENCES teams(id),
  home_score INT DEFAULT 0,
  guest_score INT DEFAULT 0,
  status TEXT DEFAULT 'SCHEDULED', -- SCHEDULED, ONGOING, COMPLETED
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

-- Create index
CREATE INDEX IF NOT EXISTS idx_fixtures_season_id ON fixtures(season_id);
CREATE INDEX IF NOT EXISTS idx_fixtures_round ON fixtures(round);

-- ==================== Matches Table ====================
CREATE TABLE IF NOT EXISTS matches (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  fixture_id UUID NOT NULL REFERENCES fixtures(id) ON DELETE CASCADE,
  match_order INT NOT NULL,
  match_type TEXT NOT NULL, -- SINGLE, DOUBLE
  home_player1_id UUID NOT NULL REFERENCES players(id),
  home_player2_id UUID REFERENCES players(id),
  guest_player1_id UUID NOT NULL REFERENCES players(id),
  guest_player2_id UUID REFERENCES players(id),
  handicap_snapshot TEXT,
  rank_snapshot TEXT,        -- JSON: player ranks at match time
  point_before TEXT,         -- JSON: points before match
  point_after TEXT,          -- JSON: points after match
  home_sets INT[],
  guest_sets INT[],
  winner_team_id UUID REFERENCES teams(id),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

-- Create index
CREATE INDEX IF NOT EXISTS idx_matches_fixture_id ON matches(fixture_id);
CREATE INDEX IF NOT EXISTS idx_matches_winner ON matches(winner_team_id);
CREATE INDEX IF NOT EXISTS idx_matches_created_at ON matches(created_at DESC);

-- ==================== Staging Players Table ====================
-- For importing raw data from Excel/CSV
CREATE TABLE IF NOT EXISTS staging_players (
  id SERIAL PRIMARY KEY,
  vdv_ten TEXT,
  nam_sinh INT,
  vdv_hang VARCHAR(10),
  diem_tich_luy NUMERIC(10,2),
  doi_bong_ten TEXT,
  mua_giai_ten TEXT,
  trang_thai_thi_dau TEXT,
  stt INT,
  created_at TIMESTAMP DEFAULT now()
);

-- ==================== Create Admin User (optional) ====================
-- This is a placeholder for future authentication
CREATE TABLE IF NOT EXISTS admins (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT DEFAULT 'ADMIN', -- ADMIN, MODERATOR, VIEWER
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

-- ==================== Useful Views ====================

-- View for getting top scorers
CREATE OR REPLACE VIEW v_season_leaderboard AS
SELECT 
  ROW_NUMBER() OVER (PARTITION BY ps.season_id ORDER BY ps.accumulated_points DESC) as rank,
  ps.id as player_season_id,
  ps.season_id,
  s.name as season_name,
  p.id as player_id,
  p.full_name,
  p.avatar_url,
  t.id as team_id,
  t.name as team_name,
  r.id as rank_id,
  r.description as rank_name,
  ps.accumulated_points,
  ps.status,
  ps.display_order
FROM player_seasons ps
JOIN seasons s ON s.id = ps.season_id
JOIN players p ON p.id = ps.player_id
JOIN teams t ON t.id = ps.team_id
JOIN ranks r ON r.id = ps.rank_id
WHERE ps.status = 'ACTIVE';

-- View for fixture details
CREATE OR REPLACE VIEW v_fixture_details AS
SELECT 
  f.id,
  f.season_id,
  f.round,
  ht.id as home_team_id,
  ht.name as home_team_name,
  gt.id as guest_team_id,
  gt.name as guest_team_name,
  f.home_score,
  f.guest_score,
  f.status,
  f.created_at
FROM fixtures f
JOIN teams ht ON ht.id = f.home_team_id
JOIN teams gt ON gt.id = f.guest_team_id;

-- ==================== Materialized View for Leaderboard (Performance Optimization) ====================
-- Use for high-frequency reads in large tournaments (300+ players)
-- Refresh every 5 minutes via background job

CREATE MATERIALIZED VIEW IF NOT EXISTS mat_season_leaderboard AS
SELECT 
  ROW_NUMBER() OVER (PARTITION BY ps.season_id ORDER BY ps.accumulated_points DESC) as rank,
  ps.id as player_season_id,
  ps.season_id,
  s.name as season_name,
  p.id as player_id,
  p.full_name as player_name,
  p.avatar_url,
  t.id as team_id,
  t.name as team_name,
  r.id as rank_id,
  r.description as rank_name,
  ps.accumulated_points,
  ps.status,
  ps.display_order,
  ps.created_at as player_season_created_at,
  NOW() as materialized_at
FROM player_seasons ps
JOIN seasons s ON s.id = ps.season_id
JOIN players p ON p.id = ps.player_id
JOIN teams t ON t.id = ps.team_id
JOIN ranks r ON r.id = ps.rank_id
WHERE ps.status = 'ACTIVE';

-- Create index on materialized view for better performance
CREATE INDEX IF NOT EXISTS idx_mat_leaderboard_season_id 
  ON mat_season_leaderboard(season_id, rank);

-- ==================== Grant Refresh Permission ====================
-- Allow background jobs to refresh the materialized view
-- GRANT REFRESH ON MATERIALIZED VIEW mat_season_leaderboard TO postgres;

