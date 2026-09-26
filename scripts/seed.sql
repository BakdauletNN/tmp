-- Demo data for local development only. This script resets the application tables.
BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
TRUNCATE bookings, desks, rooms, offices, users RESTART IDENTITY CASCADE;

INSERT INTO users (id, name, email, who, tg_chat_id, pass_hash, is_registered) VALUES
 (1, 'Admin CoworkGo', 'admin@coworkgo.kz', 'admin', 0, crypt('admin1234', gen_salt('bf', 10)), false),
 (2, 'Aigerim Demo', 'demo@coworkgo.kz', '', 0, crypt('demo1234', gen_salt('bf', 10)), false),
 (3, 'Timur Client', 'timur@coworkgo.kz', '', 0, crypt('demo1234', gen_salt('bf', 10)), false);

INSERT INTO offices (id, name, address, star, has_kitchen, time_range_work, metro_near) VALUES
 (1, 'Alatau Hub', 'Rozybakiyeva St. 247a, Almaty', 5, true, '08:00–22:00', false),
 (2, 'Abai Space', 'Abay Ave. 44, Almaty', 5, true, '24/7', true),
 (3, 'Dostyk Loft', 'Dostyk Ave. 132, Almaty', 4, true, '09:00–21:00', true),
 (4, 'Nauryz Point', 'Tole Bi St. 273, Almaty', 4, false, '08:00–20:00', false),
 (5, 'Samal Workroom', 'Samal-2 District, 111, Almaty', 5, true, '07:00–23:00', false),
 (6, 'Zhibek Zholy Desk', 'Zhibek Zholy St. 135, Almaty', 3, false, '09:00–18:00', true),
 (7, 'Al-Farabi Tower', 'Al-Farabi Ave. 77/7, Almaty', 5, true, '24/7', false),
 (8, 'Sairan Base', 'Raimbek St. 212, Almaty', 3, false, '10:00–20:00', true);

INSERT INTO rooms (id, office_id, type, price_hour, qty_desks, qty_person, access_code, has_air_conditioner, has_prayer_room) VALUES
 (1, 1, 'private', 6000, 4, 4, 4821, true, false), (2, 1, 'open', 2500, 12, 12, NULL, true, true),
 (3, 1, 'meeting', 9000, 8, 8, NULL, true, false), (4, 2, 'private', 7500, 6, 6, 1357, true, true),
 (5, 2, 'open', 3000, 16, 16, NULL, true, true), (6, 2, 'meeting', 12000, 10, 10, NULL, true, false),
 (7, 2, 'private', 5000, 2, 2, 9024, true, false), (8, 3, 'open', 2200, 10, 10, NULL, false, false),
 (9, 3, 'private', 5500, 4, 4, 2468, true, false), (10, 3, 'meeting', 8000, 6, 6, NULL, true, true),
 (11, 4, 'open', 1800, 14, 14, NULL, false, true), (12, 4, 'private', 4500, 3, 3, 7312, false, false),
 (13, 4, 'meeting', 6500, 8, 8, NULL, true, false), (14, 5, 'open', 2800, 20, 20, NULL, true, true),
 (15, 5, 'private', 8000, 6, 6, 1590, true, true), (16, 5, 'meeting', 11000, 12, 12, NULL, true, false),
 (17, 5, 'private', 4000, 2, 2, 6643, true, false), (18, 6, 'open', 1500, 8, 8, NULL, false, false),
 (19, 6, 'private', 3800, 3, 3, 5187, false, false), (20, 6, 'meeting', 5500, 6, 6, NULL, false, true),
 (21, 7, 'open', 3500, 24, 24, NULL, true, true), (22, 7, 'private', 9500, 8, 8, 3075, true, false),
 (23, 7, 'meeting', 14000, 14, 14, NULL, true, true), (24, 7, 'private', 6000, 4, 4, 8420, true, true),
 (25, 8, 'open', 1600, 10, 10, NULL, false, false), (26, 8, 'private', 3500, 2, 2, 4491, false, true),
 (27, 8, 'meeting', 5000, 6, 6, NULL, true, false);

INSERT INTO desks (room_id, label)
SELECT r.id, 'Desk ' || n
FROM rooms r CROSS JOIN LATERAL generate_series(1, r.qty_desks) AS n
WHERE r.type = 'open';

INSERT INTO bookings (id, room_id, user_id, desk_id, start_time, end_time, public_code) VALUES
 (1, 2, 2, (SELECT min(id) FROM desks WHERE room_id = 2), date_trunc('day', now()) - interval '5 days' + interval '10 hours', date_trunc('day', now()) - interval '5 days' + interval '13 hours', 'CG-D0001'),
 (2, 4, 2, NULL, date_trunc('day', now()) - interval '2 days' + interval '14 hours', date_trunc('day', now()) - interval '2 days' + interval '16 hours', 'CG-D0002'),
 (3, 6, 2, NULL, date_trunc('day', now()) + interval '1 day' + interval '15 hours', date_trunc('day', now()) + interval '1 day' + interval '17 hours', 'CG-D0003'),
 (4, 21, 2, (SELECT min(id) FROM desks WHERE room_id = 21), date_trunc('day', now()) + interval '2 days' + interval '9 hours', date_trunc('day', now()) + interval '2 days' + interval '18 hours', 'CG-D0004'),
 (5, 9, 2, NULL, date_trunc('day', now()) + interval '4 days' + interval '11 hours', date_trunc('day', now()) + interval '4 days' + interval '13 hours', 'CG-D0005'),
 (6, 1, 3, NULL, date_trunc('day', now()) + interval '1 day' + interval '10 hours', date_trunc('day', now()) + interval '1 day' + interval '12 hours', 'CG-T0001'),
 (7, 5, 3, (SELECT min(id) FROM desks WHERE room_id = 5), date_trunc('day', now()) + interval '1 day' + interval '12 hours', date_trunc('day', now()) + interval '1 day' + interval '14 hours', 'CG-T0002'),
 (8, 14, 3, (SELECT min(id) FROM desks WHERE room_id = 14), date_trunc('day', now()) + interval '3 days' + interval '9 hours', date_trunc('day', now()) + interval '3 days' + interval '12 hours', 'CG-T0003'),
 (9, 22, 3, NULL, date_trunc('day', now()) + interval '2 days' + interval '13 hours', date_trunc('day', now()) + interval '2 days' + interval '15 hours', 'CG-T0004');

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT max(id) FROM users));
SELECT setval(pg_get_serial_sequence('offices', 'id'), (SELECT max(id) FROM offices));
SELECT setval(pg_get_serial_sequence('rooms', 'id'), (SELECT max(id) FROM rooms));
SELECT setval(pg_get_serial_sequence('bookings', 'id'), (SELECT max(id) FROM bookings));
COMMIT;