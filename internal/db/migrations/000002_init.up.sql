-- Enable UUID generation (choose one)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- USERS
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username VARCHAR(32) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- DILEMMAS
CREATE TABLE IF NOT EXISTS dilemmas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (char_length(title) <= 285),
    is_anonymous BOOLEAN NOT NULL DEFAULT false,
    status TEXT NOT NULL DEFAULT 'active' CHECK (
        status IN ('active', 'removed')
    ),
    total_votes INTEGER NOT NULL DEFAULT 0 CHECK (total_votes >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_dilemmas_created_at ON dilemmas (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_dilemmas_status_created_at ON dilemmas (status, created_at DESC);

-- OPTIONS
CREATE TABLE IF NOT EXISTS dilemma_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    dilemma_id UUID NOT NULL REFERENCES dilemmas (id) ON DELETE CASCADE,
    label TEXT NOT NULL CHECK (char_length(label) <= 120),
    vote_count INTEGER NOT NULL DEFAULT 0 CHECK (vote_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_options_dilemma_id ON dilemma_options (dilemma_id);

-- VOTES (1 vote per user per dilemma)
CREATE TABLE IF NOT EXISTS votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    dilemma_id UUID NOT NULL REFERENCES dilemmas (id) ON DELETE CASCADE,
    option_id UUID NOT NULL REFERENCES dilemma_options (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, dilemma_id)
);

CREATE INDEX IF NOT EXISTS idx_votes_dilemma_id ON votes (dilemma_id);

CREATE INDEX IF NOT EXISTS idx_votes_option_id ON votes (option_id);

-- COMMENTS (optional, for future use)
-- COMMENTS (optional MVP, included)
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    dilemma_id UUID NOT NULL REFERENCES dilemmas (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content TEXT NOT NULL CHECK (char_length(content) <= 500),
    parent_id UUID NULL REFERENCES comments (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_comments_dilemma_id_created_at ON comments (dilemma_id, created_at DESC);

-- REPORTS (basic moderation)
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    reporter_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    dilemma_id UUID NULL REFERENCES dilemmas (id) ON DELETE CASCADE,
    comment_id UUID NULL REFERENCES comments (id) ON DELETE CASCADE,
    reason TEXT NOT NULL CHECK (char_length(reason) <= 200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (
        (
            dilemma_id IS NOT NULL
            AND comment_id IS NULL
        )
        OR (
            dilemma_id IS NULL
            AND comment_id IS NOT NULL
        )
    )
);

CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports (created_at DESC);