-- +goose Up
INSERT INTO rooms (number, floor, room_size, status) VALUES
('101', '1', 25.5, 'available'),
('102', '1', 30.0, 'available'),
('201', '2', 40.0, 'occupied');

INSERT INTO guests (first_name, last_name, room_id) VALUES
('John', 'Doe', 1),
('Jane', 'Smith', 2);

-- +goose Down
DELETE FROM guests;
DELETE FROM rooms;