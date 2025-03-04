-- NOTE(patrik): Final

-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL COLLATE NOCASE CHECK(username<>'') UNIQUE,
    password TEXT NOT NULL CHECK(password<>''),
    role TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE users_settings (
    id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT
);

CREATE TABLE api_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE series (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL CHECK(name<>''),

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE seasons (
    id TEXT PRIMARY KEY,
    serie_id TEXT NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    number INTEGER NOT NULL, 
    type TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE episodes (
    id TEXT PRIMARY KEY,
    season_id TEXT NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    -- TODO(patrik): Should these be nullable?
    season_number INTEGER,
    serie_number INTEGER,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE media (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,
    file_modified_time INTEGER NOT NULL,

    chapters TEXT NOT NULL,
    subtitles TEXT NOT NULL,
    attachments TEXT NOT NULL,
    video_tracks TEXT NOT NULL,
    audio_tracks TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE media_variants (
    id TEXT PRIMARY KEY,
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    language TEXT NOT NULL,
    video_track INTEGER NOT NULL,
    audio_track INTEGER NOT NULL,
    subtitle INTEGER,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE media_episodes (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    episode_id TEXT NOT NULL REFERENCES episodes(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE media_episodes;
DROP TABLE media_variants;
DROP TABLE media;

DROP TABLE episodes;
DROP TABLE seasons;
DROP TABLE series;

DROP TABLE api_tokens;
DROP TABLE users_settings;
DROP TABLE users;

