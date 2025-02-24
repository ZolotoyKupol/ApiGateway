-- +goose Up
CREATE TABLE IF NOT EXISTS rooms (
    id TEXT PRIMARY KEY,
    number TEXT NOT NULL,
    floor TEXT NOT NULL,
    room_size FLOAT NOT NULL,
    status TEXT NOT NULL,
    occupied_by TEXT,
    check_in TIMESTAMP,
    check_out TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guests (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    room_id INT REFERENCES rooms(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS guests;
DROP TABLE IF EXISTS rooms;