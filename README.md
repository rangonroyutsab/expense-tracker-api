# Personal Expense Tracker API

## Overview

Personal Expense Tracker API is a Go REST API built with Beego v2. It provides user registration, login, expense CRUD, date filtering, sorting, spending summaries, Swagger documentation, and focused automated tests.

This is an assignment-friendly implementation that intentionally keeps storage and authentication simple: data is persisted in CSV files, and protected expense routes use an `X-User-ID` request header instead of token-based auth.

## Features

- User registration and login
- CSV-backed user and expense persistence
- Expense create, read, update, and delete endpoints
- User-scoped expense access through `X-User-ID`
- Date-range filtering for expense lists and summaries
- Sorting by amount or expense date
- Category-based expense summary
- Swagger UI and generated OpenAPI files
- Unit and controller/router tests

## Tech Stack

- Go `1.25.0`
- Beego v2
- Bee CLI for local development and Swagger generation
- CSV files for persistence
- Go `testing` package and `httptest`

## Quick Start

Prerequisites:

- Go `1.25.0`
- Bee CLI available as `bee`

Install Bee if needed:

```bash
go install github.com/beego/bee/v2@latest
```

Set up and run the API:

```bash
cp conf/app.conf.sample conf/app.conf
go mod tidy
bee run -gendoc=true -downdoc=true
```

The server runs on port `8080` by default. Check it with:

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:

```json
{
  "success": true,
  "message": "Server is running"
}
```

## Configuration

Copy the sample config before running locally:

```bash
cp conf/app.conf.sample conf/app.conf
```

Default values:

```ini
appname = expense-tracker-api
httpport = 8080
runmode = dev
copyrequestbody = true

EnableDocs = true

users_csv = data/users.csv
expenses_csv = data/expenses.csv
```

`conf/app.conf` is ignored by Git so local settings can differ from the sample. The configured CSV files are created automatically at startup if they do not exist.

## Project Structure

```text
expense-tracker-api/
├── conf/
│   └── app.conf.sample
├── controllers/
│   ├── auth.go
│   ├── auth_test.go
│   ├── base.go
│   ├── dto.go
│   ├── expense.go
│   ├── expense_test.go
│   └── health.go
├── models/
│   ├── config.go
│   ├── expense.go
│   ├── expense_test.go
│   ├── user.go
│   └── user_test.go
├── routers/
│   ├── router.go
│   └── router_test.go
├── go.mod
├── go.sum
├── main.go
└── README.md
```

Generated or local-only paths such as `conf/app.conf`, `data/*.csv`, `swagger/`, the compiled `expense-tracker-api` binary, and `coverage.out` are ignored by Git.

## API Reference

All API routes use the `/api/v1` base path.

| Method | Path | Auth | Success |
| --- | --- | --- | --- |
| `GET` | `/api/v1/health` | No | `200 OK` |
| `POST` | `/api/v1/auth/register` | No | `201 Created` |
| `POST` | `/api/v1/auth/login` | No | `200 OK` |
| `POST` | `/api/v1/expenses` | Yes | `201 Created` |
| `GET` | `/api/v1/expenses` | Yes | `200 OK` |
| `GET` | `/api/v1/expenses/:id` | Yes | `200 OK` |
| `PUT` | `/api/v1/expenses/:id` | Yes | `200 OK` |
| `DELETE` | `/api/v1/expenses/:id` | Yes | `200 OK` |
| `GET` | `/api/v1/expenses/summary` | Yes | `200 OK` |

## Authentication

Registering creates a user record in `users_csv`. Login validates email and password and returns the user's ID:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user_id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

For every expense endpoint, send that ID in the `X-User-ID` header:

```text
X-User-ID: 1
```

The API checks that the header exists, is a positive integer, and matches a user in `users_csv`. There are no sessions, JWTs, or password hashes in this project.

## Request and Response Format

Use JSON request bodies with:

```text
Content-Type: application/json
```

Successful responses use:

```json
{
  "success": true,
  "message": "Expense created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Lunch",
    "amount": 350.5,
    "category": "Food",
    "note": "Team lunch",
    "expense_date": "2025-06-10",
    "created_at": "2026-06-02T00:00:00Z"
  }
}
```

Error responses use:

```json
{
  "success": false,
  "message": "Unauthorized"
}
```

Representative requests:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "title": "Lunch",
    "amount": 350.50,
    "category": "Food",
    "note": "Team lunch",
    "expense_date": "2025-06-10"
  }'
```

```bash
curl "http://localhost:8080/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30&sort_by=amount&sort_order=desc" \
  -H "X-User-ID: 1"
```

```bash
curl "http://localhost:8080/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: 1"
```

## Query Parameters

Expense list endpoint:

```text
GET /api/v1/expenses
```

| Parameter | Required | Description |
| --- | --- | --- |
| `date_from` | No | Include expenses on or after this date. Format: `YYYY-MM-DD`. |
| `date_to` | No | Include expenses on or before this date. Format: `YYYY-MM-DD`. |
| `sort_by` | No | Sort field. Allowed values: `amount`, `expense_date`. |
| `sort_order` | No | Sort direction. Allowed values: `asc`, `desc`. Defaults to `desc` when `sort_by` is present. |

Rules:

- If `sort_by` is omitted, results remain in CSV order.
- `sort_order` cannot be used without `sort_by`.
- Dates must use `YYYY-MM-DD`.

Expense summary endpoint:

```text
GET /api/v1/expenses/summary
```

| Parameter | Required | Description |
| --- | --- | --- |
| `date_from` | No | Summary start date. Format: `YYYY-MM-DD`. |
| `date_to` | No | Summary end date. Format: `YYYY-MM-DD`. |

For a date-range summary, provide both `date_from` and `date_to`. Supplying only one returns `400 Bad Request`.

Summary responses include:

```json
{
  "success": true,
  "message": "Summary generated",
  "data": {
    "date_from": "2025-06-01",
    "date_to": "2025-06-30",
    "total_amount": 400.5,
    "total_count": 2,
    "by_category": [
      {
        "category": "Food",
        "total": 350.5,
        "count": 1
      },
      {
        "category": "Transport",
        "total": 50,
        "count": 1
      }
    ]
  }
}
```

## CSV Storage

CSV paths are configured in `conf/app.conf`:

```ini
users_csv = data/users.csv
expenses_csv = data/expenses.csv
```

User CSV header:

```csv
id,name,email,password,created_at
```

Expense CSV header:

```csv
id,user_id,title,amount,category,note,expense_date,created_at
```

The application creates the containing directory and header rows on startup. Tests override these paths with temporary files for isolation.

## Swagger Documentation

Swagger generation is handled by Bee:

```bash
bee run -gendoc=true -downdoc=true
```

When the server is running in `dev` mode, generated docs are served from:

- UI: `http://localhost:8080/swagger/`
- JSON: `http://localhost:8080/swagger/swagger.json`
- YAML: `http://localhost:8080/swagger/swagger.yml`

The `swagger/` directory is generated locally and ignored by Git.

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```

Current tests cover model CSV behavior, auth controller behavior, expense controller behavior, query validation, summaries, and route registration.

## Validation and Errors

Registration validation:

- `name` is required.
- `email` is required and must parse as a valid email address.
- `password` is required and must be at least 6 characters.
- Duplicate emails return `409 Conflict`.

Login validation:

- `email` and `password` are required.
- Unknown users or wrong passwords return `401 Unauthorized`.

Expense validation:

- `title` is required.
- `amount` must be greater than zero.
- `category` is required and must be one of: `Food`, `Transport`, `Housing`, `Entertainment`, `Shopping`, `Healthcare`, `Education`, `Utilities`, `Other`.
- `expense_date` is required and must use `YYYY-MM-DD`.
- Invalid, missing, zero, or unknown `X-User-ID` values return `401 Unauthorized`.
- Invalid expense IDs return `400 Bad Request`; missing expenses return `404 Not Found`.

Common error examples:

```json
{
  "success": false,
  "message": "Category is invalid"
}
```

```json
{
  "success": false,
  "message": "Both date_from and date_to are required for date range summary"
}
```

## Development Notes

- Runtime routes are registered in `routers/router.go` with Beego namespaces.
- The `swaggerDocsNamespace` helper exists only so Bee can generate docs from controller annotations.
- `models.InitPaths()` loads CSV file paths from Beego config before files are initialized.
- CSV records use incrementing integer IDs based on the maximum existing ID.
- Amounts are written to CSV with two decimal places.

## Known Limitations

- Passwords are stored in plain text.
- `X-User-ID` is a demonstration auth mechanism and can be spoofed by any client that knows a user ID.
- There is no JWT/session handling, password reset flow, role model, or HTTPS enforcement.
- CSV persistence has no database indexes, transactions, or concurrent write protection.
- Expense listing has no pagination.
- This implementation is suitable for learning and evaluation, not production deployment.

## Future Improvements

- Hash passwords with bcrypt or Argon2.
- Replace header-based auth with JWT or server-side sessions.
- Move persistence to a database such as PostgreSQL or SQLite.
- Add pagination and richer filtering for expense lists.
- Add migration/seeding support.
- Add structured logging and request tracing.
- Add Docker and CI configuration.
