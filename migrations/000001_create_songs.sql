-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_title VARCHAR(255) NOT NULL,
    release_date VARCHAR(50),
    lyrics TEXT,
    youtube_link VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_songs_group ON songs(group_name);
CREATE INDEX idx_songs_title ON songs(song_title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd