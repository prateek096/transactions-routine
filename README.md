# transactions-routine
Service to manage accounts and transactions.

## Run (Docker)
- Start: `make run` (runs `docker compose up --build`)
- Stop: `make stop`

## Env
- `DB_URL` (e.g. `postgres://user_test:password123@db:5432/finance_db?sslmode=disable`)
- `PORT` (default: `8080`)

## Endpoints
- `POST /accounts` — create account
  - JSON: `{ "document_number": "12345678901" }`
- `GET /accounts/:accountId` — get account by id
- `POST /transactions` — create transaction
  - JSON: `{ "account_id": 1, "operation_type_id": 4, "amount": 100.00 }`


