import os
from contextlib import contextmanager
from typing import Any, Dict, Iterator, List, Optional, Tuple

import psycopg
from psycopg.rows import dict_row


BOOKING_QUERY = """
    SELECT b.public_code, b.start_time, b.end_time,
           r.type AS room_type, o.name AS office_name, o.address
    FROM bookings b
    JOIN rooms r ON r.id = b.room_id
    JOIN offices o ON o.id = r.office_id
    WHERE b.public_code = %s
"""


class Storage:
    def __init__(self, database_url: Optional[str] = None) -> None:
        self.database_url = database_url or os.environ["DB_URL"]

    def get_booking(self, code: str) -> Optional[Dict[str, Any]]:
        with psycopg.connect(self.database_url, row_factory=dict_row) as conn:
            return conn.execute(BOOKING_QUERY, (code,)).fetchone()

    def link_chat(self, chat_id: int, code: str) -> Optional[str]:
        with psycopg.connect(self.database_url, row_factory=dict_row) as conn:
            with conn.transaction():
                booking = conn.execute(
                    "SELECT user_id FROM bookings WHERE public_code = %s FOR UPDATE",
                    (code,),
                ).fetchone()
                if booking is None:
                    return None
                conn.execute(
                    "UPDATE users SET tg_chat_id = 0 WHERE tg_chat_id = %s",
                    (chat_id,),
                )
                user = conn.execute(
                    "UPDATE users SET tg_chat_id = %s WHERE id = %s RETURNING name",
                    (chat_id, booking["user_id"]),
                ).fetchone()
                return user["name"] if user else None

    def unlink_chat(self, chat_id: int) -> bool:
        with psycopg.connect(self.database_url) as conn:
            result = conn.execute(
                "UPDATE users SET tg_chat_id = 0 WHERE tg_chat_id = %s",
                (chat_id,),
            )
            return result.rowcount > 0

    @contextmanager
    def lock_due_reminders(
        self, limit: int = 50
    ) -> Iterator[Tuple[List[Dict[str, Any]], Any]]:
        with psycopg.connect(self.database_url, row_factory=dict_row) as conn:
            with conn.transaction():
                rows = conn.execute(
                    """
                    SELECT b.id, b.public_code, b.start_time, u.tg_chat_id,
                           o.name AS office_name
                    FROM bookings b
                    JOIN users u ON u.id = b.user_id
                    JOIN rooms r ON r.id = b.room_id
                    JOIN offices o ON o.id = r.office_id
                    WHERE b.reminder_sent_at IS NULL
                      AND b.start_time > now()
                      AND b.start_time <= now() + interval '24 hours'
                      AND COALESCE(u.tg_chat_id, 0) <> 0
                    ORDER BY b.start_time
                    LIMIT %s
                    FOR UPDATE OF b SKIP LOCKED
                    """,
                    (limit,),
                ).fetchall()
                yield [dict(row) for row in rows], conn