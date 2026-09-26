import unittest
from datetime import datetime

from notifier.logic import format_booking_summary, parse_booking_code, reminder_text


class BookingCodeTests(unittest.TestCase):
    def test_strips_whitespace_and_accepts_seed_code(self):
        self.assertEqual(parse_booking_code(" CG-D0001 "), "CG-D0001")

    def test_rejects_invalid_code(self):
        self.assertIsNone(parse_booking_code("not a code"))


class BookingMessageTests(unittest.TestCase):
    def test_formats_booking_and_reminder(self):
        booking = {
            "public_code": "CG-D0001",
            "office_name": "Abai Space",
            "room_type": "meeting",
            "start_time": datetime(2026, 10, 1, 10, 0),
            "end_time": datetime(2026, 10, 1, 12, 0),
        }
        summary = format_booking_summary(booking)
        reminder = reminder_text(booking)

        self.assertIn("Abai Space", summary)
        self.assertIn("01.10.2026 10:00 - 01.10.2026 12:00", summary)
        self.assertIn("CG-D0001", reminder)
        self.assertIn("is tomorrow", reminder)


if __name__ == "__main__":
    unittest.main()