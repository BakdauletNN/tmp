.PHONY: help dev api bot worker front infra migrate setup setup-go setup-python setup-frontend test test-go test-python test-frontend build build-go build-python build-frontend check

help:
	@printf '%s\n' \
		'make setup   Install Go, Python, and frontend dependencies' \
		'make infra   Start PostgreSQL and Redis' \
		'make migrate Apply database migrations' \
		'make dev     Start API, Telegram bot, worker, and frontend' \
		'make test    Run Go, Python, and frontend checks' \
		'make build   Build/check all project components' \
		'make check   Run the complete test and build sequence'

dev:
	$(MAKE) -j4 api bot worker front

api:
	go run ./cmd/api

bot:
	PYTHONPATH=cmd/worker .venv/bin/python -m notifier.bot

worker:
	PYTHONPATH=cmd/worker .venv/bin/python -m notifier.worker

front:
	cd front && npm run dev

infra:
	docker compose up -d

migrate:
	bash scripts/migrate.sh

setup: setup-go setup-python setup-frontend

setup-go:
	go mod download

setup-frontend:
	cd front && npm ci

test: test-go test-python test-frontend

test-go:
	go test ./...

test-python:
	PYTHONPATH=cmd/worker .venv/bin/python -m unittest discover -s cmd/worker/notifier -p 'test_*.py' -v

test-frontend:
	cd front && npm run build

build: build-go build-python build-frontend

build-go:
	go build ./...

build-python:
	.venv/bin/python -m compileall -q cmd/worker/notifier

build-frontend:
	cd front && npm run build

check:
	$(MAKE) test
	$(MAKE) build

setup-python:
	python3 -m venv .venv
	.venv/bin/pip install -r cmd/worker/notifier/requirements.txt