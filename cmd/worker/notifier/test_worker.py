import asyncio
import unittest
from contextlib import contextmanager

from notifier.worker import deliver_due_reminders


class FakeConnection:
    def __init__(self):
        self.statements = []

    def execute(self, query, params):
        self.statements.append((query, params))


class FakeStorage:
    def __init__(self, bookings):
        self.connection = FakeConnection()
        self.bookings = bookings

    @contextmanager
    def lock_due_reminders(self):
        yield self.bookings, self.connection


class FakeBot:
    def __init__(self, fail=False):
        self.messages = []
        self.fail = fail

    async def send_message(self, chat_id, text):
        if self.fail:
            raise RuntimeError("Telegram unavailable")
        self.messages.append((chat_id, text))


class ReminderWorkerTests(unittest.TestCase):
    def setUp(self):
        self.booking = {
            "id": 12,
            "public_code": "a1b2c3d4e5f6",
            "start_time": "2026-10-01 10:00",
            "office_name": "Abai Space",
            "tg_chat_id": 44,
        }

    def test_sends_reminder_and_marks_booking(self):
        storage = FakeStorage([self.booking])
        bot = FakeBot()
        sent = asyncio.run(deliver_due_reminders(bot, storage))

        self.assertEqual(sent, 1)
        self.assertEqual(bot.messages[0][0], 44)
        self.assertIn("Abai Space", bot.messages[0][1])
        self.assertEqual(storage.connection.statements[0][1], (12,))

    def test_failed_delivery_does_not_mark_booking(self):
        storage = FakeStorage([self.booking])
        bot = FakeBot(fail=True)

        with self.assertRaisesRegex(RuntimeError, "Telegram unavailable"):
            asyncio.run(deliver_due_reminders(bot, storage))

        self.assertEqual(storage.connection.statements, [])


if __name__ == "__main__":
    unittest.main()