ALTER TABLE bookings ADD COLUMN reminder_sent_at TIMESTAMPTZ;
CREATE INDEX bookings_due_reminders_idx
    ON bookings (start_time)
    WHERE reminder_sent_at IS NULL;