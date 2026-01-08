# Go Backend API

A production-ready Go backend API with JWT authentication, OAuth support (Google & GitHub), and PostgreSQL database.

## Features

- ✅ JWT-based authentication
- ✅ OAuth 2.0 support (Google & GitHub)
- ✅ PostgreSQL database with GORM ORM
- ✅ Database migrations with gormigrate
- ✅ Docker & Docker Compose setup
- ✅ Clean architecture with best practices
- ✅ RESTful API endpoints
- ✅ Password hashing with bcrypt
- ✅ Environment-based configuration

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/                 # HTTP handlers and routing
│   │   ├── handlers.go
│   │   └── router.go
│   ├── config/             # Configuration management
│   │   └── config.go
│   ├── database/           # Database connection and migrations
│   │   ├── database.go
│   │   └── migrations/     # Database migrations
│   │       └── migrations.go
│   ├── middleware/         # HTTP middleware
│   │   └── auth_middleware.go
│   ├── models/             # Data models (GORM)
│   │   └── user.go
│   ├── repository/         # Data access layer (GORM)
│   │   └── user_repository.go
│   └── service/            # Business logic layer
│       ├── auth_service.go
│       └── user_service.go
├── pkg/                    # Reusable packages
│   ├── auth/              # JWT authentication
│   │   └── jwt.go
│   └── oauth/             # OAuth providers
│       └── oauth.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional)
- PostgreSQL (if not using Docker)

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd studoto-backend
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start PostgreSQL Database:

**Option A: Using Docker Compose (Recommended)**

```bash
# Make sure Docker Desktop is running, then:
docker-compose up -d postgres

# Or use Make:
make docker-db

# Verify it's running:
docker-compose ps
```

**Option B: Using Local PostgreSQL**

```bash
# Install PostgreSQL (macOS)
brew install postgresql@15
brew services start postgresql@15

# Create database
createdb studoto_db
```

5. Run the application:

```bash
# Using Make:
make run

# Or directly:
go run cmd/server/main.go
```

## Environment Variables

Create a `.env` file based on `.env.example`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=backend_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24

# Server
PORT=8080

# OAuth (optional)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,X-Auth-Key
CORS_EXPOSED_HEADERS=Content-Length
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400
```

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login with email/password
- `GET /auth/oauth/:provider` - Get OAuth URL (google or github)
- `GET /auth/callback/:provider` - OAuth callback

### Protected Endpoints

- `GET /api/profile` - Get current user profile

### Example Requests

#### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Get Profile (Protected)

```bash
curl -X GET http://localhost:8080/api/profile \
  -H "X-Auth-Key: YOUR_JWT_TOKEN"
```

#### OAuth Login

```bash
# 1. Get OAuth URL
curl http://localhost:8080/auth/oauth/google

# 2. Visit the returned URL in browser
# 3. After authorization, you'll be redirected to callback with code
# 4. The callback endpoint will return JWT token
```

## Docker

### Build and Run

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Development

### Run Tests

```bash
go test ./...
```

### Build

```bash
go build -o bin/server ./cmd/server
```

### Database Migrations

Migrations are managed using gormigrate. To add a new migration:

1. Edit `internal/database/migrations/migrations.go`
2. Add a new migration to the `GetMigrations()` function:

```go
{
    ID: "20240102000001", // Use timestamp format: YYYYMMDDHHMMSS
    Migrate: func(tx *gorm.DB) error {
        // Your migration logic here
        return tx.AutoMigrate(&models.NewModel{})
    },
    Rollback: func(tx *gorm.DB) error {
        // Your rollback logic here
        return tx.Migrator().DropTable(&models.NewModel{})
    },
},
```

Migrations run automatically on application startup. The migration table tracks which migrations have been applied.

## Security Best Practices

1. **Change JWT_SECRET** in production
2. **Use strong passwords** for database
3. **Enable SSL** for database connections in production
4. **Use HTTPS** in production
5. **Validate and sanitize** all user inputs
6. **Implement rate limiting** for production
7. **Use environment variables** for sensitive data

## License

MIT
