import logging
import os
from typing import Optional

from dotenv import load_dotenv
from telegram import Update
from telegram.ext import Application, CommandHandler, ContextTypes, MessageHandler, filters

from notifier.logic import format_booking_summary, parse_booking_code
from notifier.storage import Storage


logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logging.getLogger("httpx").setLevel(logging.WARNING)
logging.getLogger("httpcore").setLevel(logging.WARNING)
logger = logging.getLogger(__name__)


def get_code_argument(context: ContextTypes.DEFAULT_TYPE) -> Optional[str]:
    return parse_booking_code(" ".join(context.args)) if context.args else None


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    if message is None:
        return
    if context.args and context.args[0].lower() == "support":
        await support(update, context)
        return
    if context.args and parse_booking_code(context.args[0]):
        await send_booking(update, context, context.args[0])
        return
    await message.reply_text(
        "Hi! Send a booking code or use a command:\n"
        "/booking CODE - booking details\n"
        "/link CODE - enable reminders\n"
        "/unlink - disable reminders\n"
        "/support - contact support"
    )


async def booking_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    code = get_code_argument(context)
    message = update.effective_message
    if code is None:
        if message is not None:
            await message.reply_text("Command format: /booking BOOKING_CODE")
        return
    await send_booking(update, context, code)


async def send_booking(update: Update, context: ContextTypes.DEFAULT_TYPE, code: str) -> None:
    message = update.effective_message
    if message is None or not message.text:
        if message is None:
            return
    booking = context.application.bot_data["storage"].get_booking(code)
    if booking is None:
        await message.reply_text("No booking found with that code.")
        return
    await message.reply_text(format_booking_summary(booking))


async def show_booking(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    if message is None or not message.text:
        return
    code = parse_booking_code(message.text)
    if code is None:
        await message.reply_text("Send a booking code or use /support.")
        return
    await send_booking(update, context, code)


async def handle_text(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    if context.chat_data.get("awaiting_support_message"):
        await forward_support_message(update, context)
        return
    await show_booking(update, context)


async def link_account(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    chat = update.effective_chat
    if message is None or chat is None:
        return
    code = get_code_argument(context)
    if code is None:
        await message.reply_text("Command format: /link BOOKING_CODE")
        return
    name = context.application.bot_data["storage"].link_chat(chat.id, code)
    if name is None:
        await message.reply_text("Could not find the booking. Check the code and try again.")
        return
    await message.reply_text(
        "All set, {name}! We will send you a reminder a day before your booking.".format(name=name)
    )


async def unlink_account(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    chat = update.effective_chat
    if message is None or chat is None:
        return
    removed = context.application.bot_data["storage"].unlink_chat(chat.id)
    await message.reply_text(
        "Reminders disabled." if removed else "This chat was not linked to an account."
    )


async def support(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    if message is None:
        return
    context.chat_data["awaiting_support_message"] = True
    await message.reply_text(
        "Write a single message describing how we can help. It will be forwarded to support."
    )


async def forward_support_message(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    chat = update.effective_chat
    if message is None or chat is None or not message.text:
        return
    support_chat_id = os.environ.get("TELEGRAM_SUPPORT_CHAT_ID", "")
    if not support_chat_id:
        context.chat_data.pop("awaiting_support_message", None)
        await message.reply_text("Support is temporarily unavailable. Please try again later.")
        logger.error("TELEGRAM_SUPPORT_CHAT_ID is not configured")
        return

    user = update.effective_user
    username = "@" + user.username if user and user.username else "no username"
    sender_name = user.full_name if user else "Customer"
    await context.bot.send_message(
        chat_id=int(support_chat_id),
        text=(
            "CoworkGo support request\n"
            "Customer: {name} ({username})\n"
            "Chat ID: {chat_id}\n\n{body}"
        ).format(name=sender_name, username=username, chat_id=chat.id, body=message.text),
    )
    context.chat_data.pop("awaiting_support_message", None)
    await message.reply_text("Your message has been forwarded to support. The reply will arrive here.")


async def reply_to_support_request(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.effective_message
    user = update.effective_user
    if message is None or user is None:
        return
    allowed_ids = {
        int(value.strip())
        for value in os.environ.get("TELEGRAM_SUPPORT_ADMIN_IDS", "").split(",")
        if value.strip().isdigit()
    }
    if user.id not in allowed_ids:
        await message.reply_text("This command is available to support operators only.")
        return
    if len(context.args) < 2 or not context.args[0].isdigit():
        await message.reply_text("Command format: /reply CHAT_ID reply text")
        return
    chat_id = int(context.args[0])
    reply = " ".join(context.args[1:])
    await context.bot.send_message(chat_id=chat_id, text="Support reply:\n" + reply)
    await message.reply_text("Reply sent.")


def main() -> None:
    load_dotenv()
    token = os.environ.get("TELEGRAM_BOT_TOKEN")
    if not token:
        raise RuntimeError("TELEGRAM_BOT_TOKEN is required")
    application = Application.builder().token(token).build()
    application.bot_data["storage"] = Storage()
    application.add_handler(CommandHandler("start", start))
    application.add_handler(CommandHandler("help", start))
    application.add_handler(CommandHandler("booking", booking_command))
    application.add_handler(CommandHandler("link", link_account))
    application.add_handler(CommandHandler("unlink", unlink_account))
    application.add_handler(CommandHandler("support", support))
    application.add_handler(CommandHandler("reply", reply_to_support_request))
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_text))
    logger.info("Telegram bot started")
    application.run_polling(allowed_updates=Update.ALL_TYPES)


if __name__ == "__main__":
    main()