# CoworkGo

A small coworking booking project with a Go backend, React frontend, and Python services.

## Requirements

* Go
* Python 3
* Node.js + npm
* Docker
* Docker Compose

## Setup

Clone the project and enter the folder:

```bash
git clone <repository-url>
cd tmp
```

Install the dependencies:

```bash
make setup
```

Create the environment file:

```bash
cp .env.template .env
```

Then fill in the required values in `.env`.

For local development, the database URL can be:

```env
DB_URL=postgres://cowork:cowork@localhost:5432/cowork
PORT=8080
REDIS_URL=redis://localhost:6379
```
## Start the project

postgres,redis:

```bash
make infra
```

Run database migrations:

```bash
make migrate
```

Then start the project:

```bash
make dev
```

The frontend should be available at:

```text
http://localhost:5173
```

The API runs on:

```text
http://localhost:8080
```

## Demo data

test data run

```bash
docker compose exec -T postgres psql -U cowork -d cowork < scripts/seed.sql
```

## Tests

Run all tests:

```bash
make test
```

Or run the complete check:

```bash
make check
```

## Build

To build/check all parts of the project:

```bash
make build
```
<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 32 54" src="https://github.com/user-attachments/assets/93590e42-5d28-4905-89cb-6c0357d4392c" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 33 02" src="https://github.com/user-attachments/assets/74ebf7c5-0422-403c-9426-9e2d059afc2f" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 34 55" src="https://github.com/user-attachments/assets/3aa0cc2e-4610-4a75-86d5-3d6aa8fde142" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 35 28" src="https://github.com/user-attachments/assets/daeb6825-48d8-4e55-b3b1-9d3b302e2397" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 36 00" src="https://github.com/user-attachments/assets/ce5a23a9-ea31-458d-b332-627e3073ac9e" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 36 13" src="https://github.com/user-attachments/assets/12f3cd0f-e18e-42d0-944f-65c42a726e9f" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 36 43" src="https://github.com/user-attachments/assets/e2ed4046-0c68-4c73-8878-c0b04f9dd98b" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 01 36 47" src="https://github.com/user-attachments/assets/769f1575-b96b-4c6a-9478-1f0f1dbf5bd3" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 02 01 45" src="https://github.com/user-attachments/assets/81cab328-4c5b-458c-a4d8-bc6747b9466d" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 02 16 09" src="https://github.com/user-attachments/assets/8baacca7-ae15-49d2-a0bd-c20188e7b08f" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 02 16 15" src="https://github.com/user-attachments/assets/5298668d-cf00-42a1-817e-960b200d60b4" />

<img width="1470" height="956" alt="Снимок экрана 2026-09-26 в 02 16 36" src="https://github.com/user-attachments/assets/b6d63782-0328-4aa8-9827-9b959136dfd8" />











