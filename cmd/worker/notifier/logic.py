import re
from datetime import datetime
from typing import Any, Dict, Optional


CODE_PATTERN = re.compile(r"^[A-Za-z0-9-]{6,32}$")


def parse_booking_code(value: str) -> Optional[str]:
    code = value.strip()
    return code if CODE_PATTERN.fullmatch(code) else None


def format_booking_summary(booking: Dict[str, Any]) -> str:
    start = _format_datetime(booking["start_time"])
    end = _format_datetime(booking["end_time"])
    return (
        "Booking found\n"
        "Code: {code}\n"
        "Space: {office}\n"
        "Room: {room}\n"
        "Time: {start} - {end}"
    ).format(
        code=booking["public_code"],
        office=booking["office_name"],
        room=booking["room_type"],
        start=start,
        end=end,
    )


def _format_datetime(value: Any) -> str:
    if isinstance(value, datetime):
        return value.strftime("%d.%m.%Y %H:%M")
    return str(value)


def reminder_text(booking: Dict[str, Any]) -> str:
    return "Reminder: your booking at {office} is tomorrow, {start}. Code: {code}".format(
        office=booking["office_name"],
        start=_format_datetime(booking["start_time"]),
        code=booking["public_code"],
    )