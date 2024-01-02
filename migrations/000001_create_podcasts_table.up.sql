CREATE TABLE IF NOT EXISTS podcasts
(
    id              BIGSERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    platform        TEXT NOT NULL,
    url             TEXT NOT NULL,
    host            TEXT NOT NULL,
    program         TEXT NOT NULL,
    guest_speakers  TEXT[] NOT NULL,
    year            INT NOT NULL,
    language        TEXT NOT NULL,
    tags            TEXT[] NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
)