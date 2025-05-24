-- +goose Up
-- Создание таблицы категорий
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Создание таблицы курсов
CREATE TABLE courses (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    level VARCHAR(50),
    duration INTEGER,
    rating FLOAT,
    students_count INTEGER DEFAULT 0,
    thumbnail VARCHAR(255),
    price DECIMAL(10, 2),
    status VARCHAR(50) DEFAULT 'draft',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Создание таблицы уроков
CREATE TABLE lessons (
    id UUID PRIMARY KEY,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    order_num INTEGER NOT NULL,
    requires_test BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (course_id, order_num)
);

-- Создание таблицы тестов
CREATE TABLE tests (
    id UUID PRIMARY KEY,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE UNIQUE,
    passing_score INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Создание таблицы вопросов
CREATE TABLE questions (
    id UUID PRIMARY KEY,
    test_id UUID REFERENCES tests(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    options JSONB,
    correct_answer INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Создание таблицы прогресса уроков
CREATE TABLE lesson_progress (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    viewed_at TIMESTAMPTZ,
    test_score INTEGER,
    passed_test BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, lesson_id)
);

-- Создание таблицы записей XP
CREATE TABLE xp_entries (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    lesson_id UUID REFERENCES lessons(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('lesson_view', 'test_pass', 'course_complete')),
    amount INTEGER NOT NULL,
    earned_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS xp_entries;
DROP TABLE IF EXISTS lesson_progress;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS categories; 