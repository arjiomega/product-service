# Product Service

A Go-based product management microservice built with Gin, PostgreSQL, and OpenAPI code generation.

Originally developed for Pabili Marketplace, this service is designed to run independently and can be integrated into other applications.

## Tech Stack

Language: Go
HTTP Framework: Gin
Database: PostgreSQL
API Specification: OpenAPI
Code Generation: oapi-codegen
Database Migrations: Goose

## Requirements

Go (version specified in go.mod)
PostgreSQL
Goose CLI


## Getting Started

### 1. Clone the repository

```bash
git clone YOUR_REPOSITORY_URL
cd pabili-product-service
```

### 2. Configure environment variables

```bash
cp .env.example .env
```

Configure your PostgreSQL connection:
```.env
DATABASE_URL=postgres://postgres:password@localhost:5432/pabili_products?sslmode=disable
```

### 3. Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Ensure that your Go binary directory is included in your system's PATH.

### 4. Run database migrations

Apply all pending migrations:
```bash
goose -dir migrations postgres \
  "postgres://postgres:password@localhost:5432/pabili_products?sslmode=disable" \
  up
```

Check migration status:
```bash
goose -dir migrations postgres \
  "postgres://postgres:password@localhost:5432/pabili_products?sslmode=disable" \
  status
```

### 5. Run the service
``` bash
go run .
```

The service starts on port 8080 by default.

## API Documentation
The OpenAPI specification is available at `openapi.yaml`

| Method | Endpoint | Description |
| - | - | - |
|GET | /health | Health check
| GET | /products | List products
| POST | /products | Create a product
| GET | /products/{product_id} | Retrieve a product

## Development

### Generate API code

API types and server interfaces are generated from the OpenAPI specification using oapi-codegen.

Run this command whenever you modify openapi.yaml:
```bash
go tool oapi-codegen -config config.yaml openapi.yaml
```

Generated code is maintained in the repository so consumers can build the service without regenerating it.

### Create a database migration

Create a new SQL migration:
```bash
goose -dir migrations create create_product_images_table sql
```

This generates a timestamped migration file in migrations/.

Add your SQL statements under the -- +goose Up and -- +goose Down sections.

### Run tests
```bash
go tests ./...
```

## License

This project is licensed under the MIT License.
