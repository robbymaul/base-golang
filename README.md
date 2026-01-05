# Application Service

A Go-based Application service with PostgreSQL database integration and Docker support.

## Prerequisites

- Go (version specified in go.mod)
- Docker and Docker Compose
- golang-migrate (will be installed automatically)
- k6 (for load testing)

## Getting Started

### Environment Setup

1. Copy the example environment file and update the values:
   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your database and service configurations.

### Installation

1. Install dependencies and configure the application:
   ```bash
   make configure
   ```

2. Download vendor dependencies:
   ```bash
   make vendor
   ```

## Development

### Running with Docker

Start the application and its dependencies:
```bash
make docker-local-up
```

Stop the application:
```bash
make docker-local-down
```

### Database Migrations

Create a new migration file:
```bash
make migrate-script
# Follow the prompt to name your migration
```

Run migrations:
```bash
# Run next migration
make migrate-up

# Run all pending migrations
make migrate-up-all

# Rollback last migration
make migrate-down

# Reset database (DANGER: drops all tables)
make migrate-clean

# Check current migration version
make migrate-version
```

### Building the Application

Build the release binary:
```bash
make release-bin
```

### Creating Admin User

Create a super admin user:
```bash
make create-admin-auth
```

## Testing

### Load Testing with k6

Run health check tests:
```bash
make k6-health-check
# When prompted, enter k6 options (e.g., "-u 10 -d 30s" for 10 users for 30 seconds)
```

Run Application scenario tests:
```bash
# Run production scenario
make k6-Application-scenario

# Run development scenario
make k6-Application-scenario-dev
```

## Project Structure

```
.
├── app/                  # Application code
├── bin/                  # Compiled binaries
├── cmd/                  # Main application entry points
├── config/               # Configuration files
├── k6/                   # Load testing scripts
│   ├── health.js         # Health check test
│   └── Application/          # Application test scenarios
│       ├── Application.js     # Production Application test
│       └── payment_dev.js # Development payment test
├── migrations/           # Database migrations
└── vendor/               # Dependencies (managed by Go modules)
```

## Development Workflow

1. Create a new migration for database changes
2. Implement your feature
3. Write tests
4. Run tests
5. Commit changes with a descriptive message
6. Push to your feature branch
7. Create a pull request

## License

[Your License Here]
