-- schema.sql (run this once to set up the database schema)

-- Artists table : the "one" side of the 
-- one-to-many relationship with albums
CREATE TABLE IF NOT EXISTS artists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Albums table : the "many" side of the
-- one-to-many relationship with artists
CREATE TABLE IF NOT EXISTS albums (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL CHECK (char_length(title) > 3),
    price NUMERIC(10, 2) NOT NULL CHECK (price > 0),
    artist_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

-- Indexes: without these, "find albums by artis"
-- scans the entire table

CREATE INDEX IF NOT EXISTS idx_albums_artist_id ON albums(artist_id) ;

-- Seed some data so the API isn't empty on first run
INSERT INTO artists (name, country) VALUES
('The Beatles', 'UK'),
('Pink Floyd', 'UK'),
('Led Zeppelin', 'UK')
-- ON CONFLICT used to avoid duplicate entries 
-- if you run this script multiple times
ON CONFLICT DO NOTHING ;

INSERT INTO albums (title, price, artist_id) VALUES
    ('Abbey Road', 56.99 , 1),
    ('Animals', 67.32 , 2),
    ('Four', 45.89 , 3)
ON CONFLICT DO NOTHING ;

