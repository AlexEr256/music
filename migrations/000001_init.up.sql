CREATE TABLE IF NOT EXISTS songs (
    song_group VARCHAR(100) NOT NULL,
    song VARCHAR(100) PRIMARY KEY, 
    song_text TEXT NOT NULL,
    link VARCHAR(200) NOT NULL,
    release_date TIMESTAMP NOT NULL
);