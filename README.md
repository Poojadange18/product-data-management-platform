# ProductHub

ProductHub is a portfolio-ready product data management application built with Go, Gin, PostgreSQL, and vanilla HTML/CSS/JavaScript.

## Features

- JWT registration and login with bcrypt password hashing.
- Protected product CRUD with manually supplied, unique product IDs.
- Validation, consistent error responses, pagination, search, category/stock filters, and safe sorting.
- Responsive dashboard, inventory health, category analytics, confirmation dialogs, loading/error feedback, and pagination controls.
- OpenAPI specification at `/openapi.yaml` and browser documentation at `/docs`.

## Run locally

1. Copy `.env.example` to `.env` and set `DB_PASSWORD` and a long random `JWT_SECRET`.
2. Create the database and run `migrations/001_initial_schema.sql` once.
3. Export the variables from `.env` in your shell, then run:

```bash
go run .
```

Open `http://localhost:8080/register` and create an account.

## Docker

Copy `.env.example` to `.env`, set secure values, then run:

```bash
docker compose up --build
```

PostgreSQL initializes the schema from `migrations/001_initial_schema.sql`. To reset the Docker database during development, use `docker compose down -v` (this removes local database data).

## API

- `GET /health` (public)
- `POST /auth/register` (public)
- `POST /auth/login` (public)
- `GET|POST /products`, `GET|PUT|DELETE /products/:id` (JWT required)

Send `Authorization: Bearer <token>` for product endpoints. `GET /products` without query parameters returns the original product-array response. Supply `page`, `limit`, `search`, `category`, `stock`, `sort`, or `order` for the paginated response. Example: `/products?page=1&limit=10&search=phone&category=Electronics&sort=price&order=asc`.

## Verify

```bash
go fmt ./...
go test ./...
go vet ./...
go build ./...
```

## Deployment

Deploy the Docker image to Render, Railway, or Fly.io with a managed PostgreSQL instance. Configure `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE=require`, `JWT_SECRET`, `SERVER_PORT`, and `CORS_ALLOWED_ORIGIN` in the platform environment settings; run the migration against the managed database before serving traffic.
