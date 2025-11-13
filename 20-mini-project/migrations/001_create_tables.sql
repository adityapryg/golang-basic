-- Migration: Create Users and Todos Tables
-- Description: Initial database schema for Todo REST API
-- Version: 001
-- Date: 2025-11-13

-- ============================================
-- CREATE USERS TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for users table
CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- ============================================
-- CREATE TODOS TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS todos (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority INTEGER NOT NULL DEFAULT 0,
    due_date TIMESTAMP NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes for todos table
CREATE INDEX idx_todos_user_id ON todos (user_id);

CREATE INDEX idx_todos_status ON todos (status);

CREATE INDEX idx_todos_deleted_at ON todos (deleted_at);

-- ============================================
-- ADD CONSTRAINTS
-- ============================================

-- Check constraint for todo status
ALTER TABLE todos
ADD CONSTRAINT check_todo_status CHECK (
    status IN (
        'pending',
        'in_progress',
        'completed'
    )
);

-- Check constraint for todo priority
ALTER TABLE todos
ADD CONSTRAINT check_todo_priority CHECK (
    priority >= 0
    AND priority <= 5
);

COMMENT ON TABLE users IS 'Stores user account information';

COMMENT ON TABLE todos IS 'Stores todo items with status tracking';

COMMENT ON COLUMN users.password IS 'Bcrypt hashed password';

COMMENT ON COLUMN todos.status IS 'Todo status: pending, in_progress, or completed';

COMMENT ON COLUMN todos.priority IS 'Priority level from 0 (lowest) to 5 (highest)';