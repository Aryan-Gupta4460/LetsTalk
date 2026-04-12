CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    username TEXT,
    content TEXT,
    room TEXT,
    timestamp TIMESTAMP
);