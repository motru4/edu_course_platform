-- +goose Up
CREATE TABLE clicker_stats (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    total_clicks BIGINT NOT NULL DEFAULT 0,
    clicks_per_second FLOAT,
    last_click_time TIMESTAMPTZ,
    last_save_time TIMESTAMPTZ,
    last_save_count INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id)
);

CREATE TABLE clicker_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    click_count INTEGER NOT NULL,
    average_cps FLOAT,
    max_cps FLOAT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE clicker_leaderboard (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255),
    score BIGINT NOT NULL,
    rank INTEGER,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id)
);

-- +goose Down
DROP TABLE IF EXISTS clicker_leaderboard;
DROP TABLE IF EXISTS clicker_sessions;
DROP TABLE IF EXISTS clicker_stats; 