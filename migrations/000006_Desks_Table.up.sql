CREATE TABLE desks (
    id SERIAL PRIMARY KEY,
    room_id INT NOT NULL REFERENCES rooms(id),
    label TEXT NOT NULL
);

ALTER TABLE bookings ADD COLUMN desk_id INT REFERENCES desks(id);