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

CREATE TABLE collections (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL CHECK(name<>''),

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE media (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,

    collection_id TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,

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

