# Chirpy

A lightweight, robust, and highly performant Twitter/X clone backend API written in Go.

---

## What This Project Does

Chirpy is a backend REST API that powers a simple social media platform. It includes the following features:
* **User Management**: Support for creating users, updating profiles, and handling password hashing.
* **Authentication**: Stateless JWT-based session management with short-lived access tokens and persistent, database-backed refresh tokens (supporting session revocation/logout).
* **Chirping**: Creating, retrieving (filtering by author, sorting chronologically), and deleting chirps (short microblog posts).
* **Webhook Integrations**: Upgrade users to a premium tier ("Chirpy Red") via external webhook integrations (e.g. Polka payment platform).
* **Static Assets Server**: Built-in file server to serve the frontend client directly.
* **Admin Dashboard & Health Checks**: Metrics dashboard tracking total page visits and internal health checks.

---

## Why You Should Care

* **Built with Go Standard Library**: Leverages Go's standard library `net/http` router (using Go 1.22+ method-based pattern routing) rather than third-party router frameworks.
* **Type-Safe DB Layer**: Uses `sqlc` to generate fully type-safe, compile-time-checked Go code directly from your database queries and schema definition.
* **Database Migrations**: Utilizes `goose` for robust database versioning and database schema migrations.
* **Robust Security Design**: Implements industry best practices for security including Argon2/bcrypt-based password hashing, short-lived JWT access tokens, and separate database-validated refresh tokens.

---

## Getting Started

### Prerequisites

To run this project locally, you will need:
* **Go** (version 1.22 or higher)
* **PostgreSQL** database
* **Goose** (database migration tool)
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```
* **SQLC** (Optional, for SQL-to-Go generation)
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

### Installation & Run Steps

1. **Clone the Repository** and navigate to the project directory:
   ```bash
   cd chirpy
   ```

2. **Configure Environment Variables**:
   Create a `.env` file in the root directory and populate it with your environment config:
   ```env
   DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
   PLATFORM="dev"
   SECRET_TOKEN="your-secure-64-character-jwt-signing-secret"
   POLKA_KEY="your-polka-payment-webhook-api-key"
   ```

3. **Run Database Migrations**:
   Run the Goose migrations to set up your PostgreSQL schema:
   ```bash
   goose -dir sql/schema postgres "$DB_URL" up
   ```

4. **Build and Start the Application**:
   Compile and run the server using:
   ```bash
   go run .
   ```
   Or build the executable first:
   ```bash
   go build -o chirpy
   ./chirpy
   ```
   The application server will start listening at `http://localhost:8080`.

---

## API Documentation

For a detailed list of all endpoints, authentication schemes, JSON payloads, and response codes, please refer to the **[API Documentation](./docs/API.md)**.