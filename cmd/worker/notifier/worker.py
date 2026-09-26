from __future__ import annotations

import asyncio
import logging
import os
from typing import TYPE_CHECKING

from notifier.logic import reminder_text

if TYPE_CHECKING:
    from telegram import Bot
    from notifier.storage import Storage


logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logging.getLogger("httpx").setLevel(logging.WARNING)
logging.getLogger("httpcore").setLevel(logging.WARNING)
logger = logging.getLogger(__name__)


async def deliver_due_reminders(bot: Bot, storage: Storage) -> int:
    delivered = 0
    with storage.lock_due_reminders() as (bookings, conn):
        for booking in bookings:
            await bot.send_message(
                chat_id=booking["tg_chat_id"],
                text=reminder_text(booking),
            )
            conn.execute(
                "UPDATE bookings SET reminder_sent_at = now() WHERE id = %s",
                (booking["id"],),
            )
            delivered += 1
    return delivered


async def run_worker() -> None:
    from dotenv import load_dotenv
    from telegram import Bot
    from notifier.storage import Storage

    load_dotenv()
    token = os.environ.get("TELEGRAM_BOT_TOKEN")
    if not token:
        raise RuntimeError("TELEGRAM_BOT_TOKEN is required")
    storage = Storage()
    async with Bot(token) as bot:
        while True:
            try:
                count = await deliver_due_reminders(bot, storage)
                if count:
                    logger.info("Sent %s booking reminder(s)", count)
            except Exception:
                logger.exception("Reminder delivery failed")
            await asyncio.sleep(int(os.environ.get("REMINDER_POLL_SECONDS", "60")))


def main() -> None:
    asyncio.run(run_worker())


if __name__ == "__main__":
    main()