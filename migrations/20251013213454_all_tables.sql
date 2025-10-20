-- +goose Up
-- +goose StatementBegin
-- ============================
-- USERS TABLE
-- ============================

CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(15) UNIQUE,
    password_hash TEXT NOT NULL,            
    avatar_url TEXT,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- ============================
-- CHATROOMS TABLE
-- ============================
CREATE TABLE chatrooms (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_group BOOLEAN,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- ============================
-- CHATROOM MEMBERS TABLE
-- ============================
CREATE TABLE chatroom_members (
    id UUID PRIMARY KEY,
    chatroom_id UUID NOT NULL REFERENCES chatrooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(chatroom_id, user_id)
);

-- ============================
-- MESSAGES TABLE
-- ============================
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    chatroom_id UUID NOT NULL REFERENCES chatrooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE INDEX idx_messages_chatroom_id ON messages(chatroom_id);
CREATE INDEX idx_messages_user_id ON messages(user_id);
-- ============================
-- SUMMARIES TABLE
-- ============================
CREATE TABLE summaries (
    id UUID PRIMARY KEY,
    chatroom_id UUID NOT NULL REFERENCES chatrooms(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    generated_by VARCHAR(50) DEFAULT 'AI',
    generated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- ============================
-- SESSIONS TABLE (Optional, JWT Refresh)
-- ============================
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS summaries;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chatroom_members;
DROP TABLE IF EXISTS chatrooms;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
