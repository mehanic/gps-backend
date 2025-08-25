CREATE TABLE tracks (
    id SERIAL PRIMARY KEY,
    device_id TEXT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    speed DOUBLE PRECISION,
    timestamp TIMESTAMP NOT NULL
);
