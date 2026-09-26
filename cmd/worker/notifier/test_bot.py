import asyncio
import logging
import unittest
from types import SimpleNamespace
from unittest.mock import patch

from notifier.bot import forward_support_message, reply_to_support_request, start


class FakeMessage:
    def __init__(self, text=""):
        self.text = text
        self.replies = []

    async def reply_text(self, text):
        self.replies.append(text)


class FakeBot:
    def __init__(self):
        self.messages = []

    async def send_message(self, chat_id, text):
        self.messages.append((chat_id, text))


class FakeStorage:
    def get_booking(self, code):
        return {
            "public_code": code,
            "office_name": "Abai Space",
            "room_type": "meeting",
            "start_time": "01.10.2026 10:00",
            "end_time": "01.10.2026 12:00",
        }


class TelegramBotTests(unittest.TestCase):
    def test_http_client_logs_do_not_include_request_urls(self):
        self.assertGreaterEqual(logging.getLogger("httpx").level, logging.WARNING)
        self.assertGreaterEqual(logging.getLogger("httpcore").level, logging.WARNING)

    def test_start_deep_link_returns_booking_details(self):
        message = FakeMessage("/start CG-D0001")
        update = SimpleNamespace(effective_message=message)
        context = SimpleNamespace(args=["CG-D0001"], application=SimpleNamespace(bot_data={"storage": FakeStorage()}))

        asyncio.run(start(update, context))

        self.assertIn("CG-D0001", message.replies[0])
        self.assertIn("Abai Space", message.replies[0])

    def test_support_message_is_forwarded_with_chat_id(self):
        message = FakeMessage("Can't find the entrance to the room")
        update = SimpleNamespace(
            effective_message=message,
            effective_chat=SimpleNamespace(id=123),
            effective_user=SimpleNamespace(username="client", full_name="Demo Client"),
        )
        bot = FakeBot()
        context = SimpleNamespace(chat_data={"awaiting_support_message": True}, bot=bot)

        with patch.dict("os.environ", {"TELEGRAM_SUPPORT_CHAT_ID": "-100456"}):
            asyncio.run(forward_support_message(update, context))

        self.assertEqual(bot.messages[0][0], -100456)
        self.assertIn("Chat ID: 123", bot.messages[0][1])
        self.assertIn("Can't find the entrance", bot.messages[0][1])
        self.assertNotIn("awaiting_support_message", context.chat_data)

    def test_only_allowlisted_operator_can_reply(self):
        message = FakeMessage("/reply 123 Reply ready")
        update = SimpleNamespace(effective_message=message, effective_user=SimpleNamespace(id=77))
        bot = FakeBot()
        context = SimpleNamespace(args=["123", "Reply", "ready"], bot=bot)

        with patch.dict("os.environ", {"TELEGRAM_SUPPORT_ADMIN_IDS": "77,88"}):
            asyncio.run(reply_to_support_request(update, context))

        self.assertEqual(bot.messages, [(123, "Support reply:\nReply ready")])


if __name__ == "__main__":
    unittest.main()