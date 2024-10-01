CREATE TABLE IF NOT EXISTS songs (
    id BIGINT PRIMARY KEY,
    song_group VARCHAR(100) NOT NULL,
    song VARCHAR(100) NOT NULL, 
    song_text TEXT NOT NULL,
    link VARCHAR(200) NOT NULL,
    release_date TIMESTAMP NOT NULL
);