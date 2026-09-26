# CoworkGo

A coworking booking service: a Go REST API, a React + TypeScript web client, a Python Telegram bot, and a Python worker for reminders.

## Running locally

1. Copy `.env.template` to `.env` and fill in `DB_URL=postgres://cowork:cowork@localhost:5432/cowork`, `PORT=8080`, `JWT_SECRET`, and `TELEGRAM_BOT_TOKEN`.
2. Install all dependencies with a single command: `make setup`.
3. Start PostgreSQL and Redis: `make infra`; apply migrations: `make migrate`.
4. For demo data, run `docker compose exec -T postgres psql -U cowork -d cowork < scripts/seed.sql`.
5. Start the API, Telegram bot, worker, and frontend with a single command: `make dev`.
6. Open `http://localhost:5173`.

The seed data contains 3 users, 8 coworking spaces, 27 rooms, 114 desks, and 9 bookings. Sign in as admin with `admin@coworkgo.kz` / `admin1234`; regular accounts are `demo@coworkgo.kz` / `demo1234` and `timur@coworkgo.kz` / `demo1234` (for testing booking conflicts). To load the demo data in pgAdmin 4, open the Query Tool for the `cowork` database and execute `scripts/seed.sql` after applying migrations. The seed script truncates application tables before inserting demo data.

In `.env`, set `VITE_TELEGRAM_BOT_USERNAME`, `TELEGRAM_SUPPORT_CHAT_ID`, and `TELEGRAM_SUPPORT_ADMIN_IDS` for the Telegram link and forwarding support requests. Bot commands: `/start CODE` and `/booking CODE` show booking details, `/link CODE` enables reminders, `/unlink` disables them, `/support` forwards a message to an operator. The operator replies with `/reply CHAT_ID text`; their Telegram ID must be included in `TELEGRAM_SUPPORT_ADMIN_IDS`. The worker checks for bookings starting within 24 hours; the bot, worker, and Telegram storage Python code live in `cmd/worker/notifier`.

For password recovery, set `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM`, and `WEB_BASE_URL` in `.env`.

### Password recovery flow

1. The client sends an email to `POST /forgot-password`. The response is the same for an existing and an unknown account, so registered addresses aren't revealed.
2. For an existing user, a random token is generated; only its SHA-256 hash is stored in the DB, valid for 30 minutes. The user's previous reset tokens are deleted.
3. SMTP sends a link to `/#/reset-password?token=...`.
4. The client sends the token and new password to `POST /reset-password`. The API checks the expiry, hashes the password with bcrypt, updates it, and deletes the token in a single transaction. The link is single-use.

### Telegram flow

`/start` shows the available commands. A booking's text code is looked up in PostgreSQL and returns a short summary. `/link CODE` finds the booking's owner and saves the Telegram chat ID on their account; `/unlink` removes the link. Once a minute, the worker selects bookings starting within the next 24 hours that haven't been notified yet, sends a message, and only marks the booking as notified once the message has been sent successfully. The bot, storage, and worker Python code lives in `cmd/worker/notifier`.

Test the whole project with a single command: `make test` (Go, Python, and frontend). Build all components: `make build`.
