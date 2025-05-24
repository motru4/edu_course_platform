-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('student','admin')),
    confirmed BOOLEAN DEFAULT false,
    google_id TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE refresh_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE TABLE purchased_courses (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID NOT NULL,
    purchased_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, course_id)
);

-- +goose Down
DROP TABLE IF EXISTS purchased_courses;
DROP TABLE IF EXISTS refresh_sessions;
DROP TABLE IF EXISTS users; 