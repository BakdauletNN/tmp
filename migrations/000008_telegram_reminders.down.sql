DROP INDEX IF EXISTS bookings_due_reminders_idx;
ALTER TABLE bookings DROP COLUMN IF EXISTS reminder_sent_at;