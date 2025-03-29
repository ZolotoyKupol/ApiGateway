-- +goose Up
-- Добавляем комнаты, если они еще не существуют
INSERT INTO rooms (number, floor, room_size, status)
SELECT '101', '1', 25.5, 'available'
WHERE NOT EXISTS (SELECT 1 FROM rooms WHERE number = '101');

INSERT INTO rooms (number, floor, room_size, status)
SELECT '102', '1', 30.0, 'available'
WHERE NOT EXISTS (SELECT 1 FROM rooms WHERE number = '102');

INSERT INTO rooms (number, floor, room_size, status)
SELECT '201', '2', 40.0, 'occupied'
WHERE NOT EXISTS (SELECT 1 FROM rooms WHERE number = '201');

-- Добавляем гостей, если они еще не существуют
INSERT INTO guests (first_name, last_name, room_id)
SELECT 'John', 'Doe', 1
WHERE NOT EXISTS (SELECT 1 FROM guests WHERE first_name = 'John' AND last_name = 'Doe');

INSERT INTO guests (first_name, last_name, room_id)
SELECT 'Jane', 'Smith', 2
WHERE NOT EXISTS (SELECT 1 FROM guests WHERE first_name = 'Jane' AND last_name = 'Smith');

-- +goose Down
-- Удаляем все записи из таблиц
DELETE FROM guests;
DELETE FROM rooms;