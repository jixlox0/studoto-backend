# Go Backend API

A production-ready Go backend API with JWT authentication, OAuth support (Google & GitHub), and PostgreSQL database.

## Features

- ✅ JWT-based authentication
- ✅ OAuth 2.0 support (Google & GitHub)
- ✅ PostgreSQL database with GORM ORM
- ✅ Database migrations with gormigrate
- ✅ Dependency injection with Google Wire
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
│       ├── main.go          # Application entry point
│       ├── wire.go          # Wire dependency injection providers
│       └── wire_gen.go      # Generated Wire code (auto-generated)
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
DB_NAME=studoto_db
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

### Web Interfaces

The project includes web-based management interfaces for PostgreSQL and Redis:

#### pgAdmin (PostgreSQL Management)
- **URL**: http://localhost:5050
- **Email**: admin@admin.com
- **Password**: admin

**PostgreSQL server is automatically configured!** The database server "Studoto PostgreSQL" should appear automatically after pgAdmin starts. If it doesn't appear:

1. Login to pgAdmin at http://localhost:5050
2. Wait a few seconds for the initialization script to complete
3. Refresh the browser or expand "Servers" in the left sidebar
4. You should see "Studoto PostgreSQL" server

**If you need to add it manually:**
1. Right-click "Servers" → "Register" → "Server"
2. **General tab:** Name = `Studoto PostgreSQL`
3. **Connection tab:**
   - Host name/address: `172.20.0.10` (PostgreSQL container IP address)
   - Port: `5432`
   - Maintenance database: `postgres`
   - Username: `postgres`
   - Password: `postgres`
4. Click "Save"

**Important:** Use `172.20.0.10` as the hostname (PostgreSQL container IP address).

#### RedisInsight (Redis Management)
- **URL**: http://localhost:8081
- Modern Redis web UI from Redis Labs

**To connect to Redis:**
1. Open http://localhost:8081 in your browser
2. Wait for RedisInsight to load (may take 10-30 seconds on first start)
3. Click "Add Redis Database" or "I already have a database"
4. Enter connection details:
   - **Host**: `172.20.0.20` (Redis container IP address)
   - **Port**: `6379`
   - **Database Alias**: `Studoto Redis` (optional)
   - Leave other fields as default
5. Click "Add Redis Database"

**Important:** 
- Use `172.20.0.20` as the hostname (Redis container IP address)
- If the page doesn't load, wait a bit longer as RedisInsight takes time to initialize
- Check container logs: `docker logs studoto_redisinsight` if issues persist

**Start web interfaces:**
```bash
# Start both interfaces
make docker-ui

# Or start individually
make docker-pgadmin      # Start pgAdmin only
make docker-redis-ui     # Start RedisInsight only
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

### Dependency Injection with Wire

This project uses [Google Wire](https://github.com/google/wire) for compile-time dependency injection.

**Generate Wire code:**

```bash
# Using Make:
make wire

# Or directly:
cd cmd/server && wire
```

**Files:**

- `cmd/server/wire.go` - Wire provider definitions (edit this to add/modify dependencies)
- `cmd/server/wire_gen.go` - Generated Wire code (auto-generated, do not edit manually)

**How it works:**

1. Wire analyzes `wire.go` to understand the dependency graph
2. Generates `wire_gen.go` with the actual dependency injection code
3. `main.go` calls `InitializeApp()` which is generated by Wire

**Adding new dependencies:**

1. Add provider functions to `wire.go`
2. Include them in the `wire.Build()` call
3. Run `make wire` to regenerate `wire_gen.go`

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
