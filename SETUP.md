# Docker Compose Setup

## Network Configuration

The docker-compose.yml uses an **external network** named `studoto` with subnet `172.28.5.0/24`.

### Create the Network

Before starting services, create the external network:

```bash
# Option 1: Use the setup script
./setup-network.sh

# Option 2: Use Makefile
make docker-network

# Option 3: Manual creation
docker network create --driver bridge --subnet=172.28.5.0/24 studoto
```

## IP Address Assignments

| Service | IP Address | Ports |
|---------|------------|-------|
| **db** (PostgreSQL) | 172.28.5.2 | 5432 |
| **cache** (Redis) | 172.28.5.4 | 6379 |
| **redisinsight** | 172.28.5.5 | - |
| **pgadmin4** | 172.28.5.6 | 5050 |
| **api** | 172.28.5.7 | 8080 |

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Database Configuration
DB_HOST=172.28.5.2
DB_PORT=5432
DB_USER=studoto
DB_PASSWORD=studoto
DB_NAME=studoto
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=172.28.5.4
REDIS_PORT=6379
REDIS_PASSWORD=Da3ZqphucUo4zw9b
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24

# Server Configuration
PORT=8080

# OAuth Configuration (optional)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,X-Auth-Key,x-auth-token,X-Auth-Token
```

## Starting Services

```bash
# 1. Create network (first time only)
make docker-network

# 2. Create .env file (if not exists)
cp .env.example .env  # Edit as needed

# 3. Start all services
docker-compose up -d

# 4. View logs
docker-compose logs -f
```

## Service Details

### PostgreSQL (db)
- **Image**: postgres:14.1-alpine
- **Database**: `studoto` (created by init.dev.sql)
- **User**: `studoto` / `studoto`
- **Healthcheck**: Checks PostgreSQL readiness

### Redis (cache)
- **Image**: redis:7.0-alpine
- **Password**: `Da3ZqphucUo4zw9b`
- **Persistence**: Saves every 10 seconds if at least 1 key changed

### RedisInsight
- **Image**: redislabs/redisinsight:latest
- **Access**: http://localhost:8081 (if port mapped)
- **Connect to Redis**: Use IP 172.28.5.4 with password `Da3ZqphucUo4zw9b`

### pgAdmin4
- **Image**: dpage/pgadmin4:8.1
- **Access**: http://localhost:5050 (if port mapped)
- **Credentials**: admin@admin.com / admin
- **Pre-configured**: Server "Studoto PostgreSQL" should appear automatically

### API
- **Build**: Uses Dockerfile with target `studoto-api`
- **Port**: 8080
- **Environment**: Loads from `.env` file
- **Dependencies**: Waits for db and cache to be healthy

## Notes

- The network is **external**, meaning it must exist before starting services
- All services use static IP addresses for predictable networking
- Redis requires password authentication (configured in docker-compose.yml)
- PostgreSQL database and user are created automatically via init.dev.sql
