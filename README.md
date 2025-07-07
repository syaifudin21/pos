# POS Application

This is a Point of Sale (POS) application built with Go.

## Getting Started

### Prerequisites

- Go (version 1.22 or higher)
- PostgreSQL database

### Environment Variables

Create a `.env` file in the root directory of the project based on `.env.example` and fill in the required environment variables, especially for database connection.

```
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_NAME=your_db_name
GOOGLE_REDIRECT_URL=your_google_redirect_url
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
JWT_SECRET=your_jwt_secret
HOST=your_application_host
```

### Running the Application

To run the API server:

```bash
go run main.go api
```

### Database Operations

#### Migrations

To run database migrations (creates and updates tables):

```bash
go run main.go migrate
```

#### Seeding

To seed the database with initial data (e.g., default payment methods, super admin user):

```bash
go run main.go seed
```

#### Reset Database

To reset the database (drops all tables, runs migrations, and seeds initial data):

```bash
go run main.go resetdb
```

## Project Structure

- `cmd/api`: Contains the main entry point for the API server.
- `cmd/migrate`: Contains the main entry point for database migrations.
- `cmd/seed`: Contains the main entry point for database seeding.
- `internal/database`: Database connection and initialization.
- `internal/handlers`: HTTP request handlers.
- `internal/middleware`: Custom Echo middleware.
- `internal/models`: Database models (structs).
- `internal/services`: Business logic and service layer.
- `internal/validators`: Request payload validation logic.
- `pkg/casbin`: Casbin authorization setup.
- `pkg/localization`: Localization utilities.
- `pkg/utils`: General utility functions (e.g., JWT, password hashing).
