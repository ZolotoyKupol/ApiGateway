CREATE TABLE IF NOT EXISTS guests (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    room_id INT NOT NULL
);