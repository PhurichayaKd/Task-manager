.PHONY: up down logs rebuild migrate psql

up:
\tdocker compose up -d --build

down:
\tdocker compose down

logs:
\tdocker compose logs -f app

rebuild:
\tdocker compose build --no-cache --pull app

migrate:
\tdocker compose exec -T db psql -U app -d taskdb < internal/db/migrate/0001_init.sql

psql:
\tdocker compose exec db psql -U app -d taskdb
