# Prime Trade Microservices

A production-ready microservices-based backend built with Go, Gin, and PostgreSQL.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
│              (Web, Mobile, Third-party Services)                │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API GATEWAY                                 │
│                      (Port 8080)                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │    CORS     │  │Rate Limiting│  │     Request Routing     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────┬───────────────────────────────────────┘
                          │
          ┌───────────────┴───────────────┐
          │                               │
          ▼                               ▼
┌─────────────────────┐       ┌─────────────────────┐
│    AUTH SERVICE     │       │    TASK SERVICE     │
│     (Port 8081)     │◄──────│     (Port 8082)     │
│                     │ Token │                     │
│  ┌───────────────┐  │Validate┌───────────────┐   │
│  │  PostgreSQL   │  │       │  PostgreSQL   │   │
│  │   (auth_db)   │  │       │   (task_db)   │   │
│  └───────────────┘  │       └───────────────┘   │
└─────────────────────┘       └─────────────────────┘
          │                               │
          └───────────────┬───────────────┘
                          │
                          ▼
              ┌─────────────────────┐
              │       REDIS         │
              │  (Caching Layer)    │
              └─────────────────────┘
```

## Project Structure

```
prime-trade-microservices/
├── api-gateway/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── middleware/
│   │   ├── proxy/
│   │   └── routes/
│   ├── Dockerfile
│   └── go.mod
├── auth-service/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── controllers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── routes/
│   │   ├── services/
│   │   └── utils/
│   ├── docs/
│   ├── Dockerfile
│   └── go.mod
├── task-service/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── controllers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── routes/
│   │   ├── services/
│   │   └── utils/
│   ├── docs/
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml
└── README.md
```

## Services

### 1. API Gateway (Port 8080)
- Single entry point for all client requests
- Request routing to appropriate microservices
- CORS handling
- Rate limiting
- Request logging
- JWT token forwarding

### 2. Auth Service (Port 8081)
- User registration and login
- JWT token generation (15-minute expiry)
- Password hashing with bcrypt
- Role management (USER, ADMIN)
- Token validation endpoint

### 3. Task Service (Port 8082)
- CRUD operations for tasks
- Authorization via Auth Service token validation
- Role-based access control
- Task status management (PENDING, IN_PROGRESS, DONE)

## API Endpoints

### Auth Service
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | User login |
| GET | `/api/v1/auth/validate` | Validate JWT token |
| GET | `/health` | Health check |

### Task Service
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tasks` | Create task |
| GET | `/api/v1/tasks` | List tasks |
| GET | `/api/v1/tasks/:id` | Get task by ID |
| PUT | `/api/v1/tasks/:id` | Update task |
| DELETE | `/api/v1/tasks/:id` | Delete task |
| GET | `/health` | Health check |

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.22+ (for local development)

### Running with Docker Compose

```bash
# Clone the repository
cd prime-trade-microservices

# Start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Local Development

```bash
# Start databases and Redis
docker-compose up postgres-auth postgres-task redis -d

# Run Auth Service
cd auth-service
go mod download
go run cmd/server/main.go

# Run Task Service (in another terminal)
cd task-service
go mod download
go run cmd/server/main.go

# Run API Gateway (in another terminal)
cd api-gateway
go mod download
go run cmd/server/main.go
```

## Environment Variables

### API Gateway
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| AUTH_SERVICE_URL | Auth service URL | http://localhost:8081 |
| TASK_SERVICE_URL | Task service URL | http://localhost:8082 |
| REDIS_URL | Redis connection URL | localhost:6379 |
| RATE_LIMIT | Requests per window | 100 |
| RATE_LIMIT_WINDOW | Window in seconds | 60 |

### Auth Service
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8081 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | Database user | auth_user |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | auth_db |
| JWT_SECRET | JWT signing secret | - |
| JWT_EXPIRY | Token expiry duration | 15m |

### Task Service
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8082 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | Database user | task_user |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | task_db |
| AUTH_SERVICE_URL | Auth service URL | http://localhost:8081 |

## API Usage Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

### Create Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "Complete documentation",
    "description": "Write comprehensive API documentation",
    "status": "PENDING"
  }'
```

### List Tasks
```bash
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Scalability Design

### Horizontal Scaling

Each microservice can be scaled independently based on load:

```yaml
# Scale specific service
docker-compose up -d --scale auth-service=3 --scale task-service=5
```

### Load Balancing

For production, use a load balancer (nginx, HAProxy, or cloud LB) in front of scaled services:

```
                    ┌─────────────────┐
                    │  Load Balancer  │
                    └────────┬────────┘
           ┌─────────────────┼─────────────────┐
           ▼                 ▼                 ▼
    ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
    │ API Gateway │   │ API Gateway │   │ API Gateway │
    │  Instance 1 │   │  Instance 2 │   │  Instance 3 │
    └─────────────┘   └─────────────┘   └─────────────┘
```

### Redis Caching Strategy

- **Session Caching**: Store validated tokens to reduce Auth Service calls
- **Response Caching**: Cache frequently accessed task lists
- **Rate Limiting**: Distributed rate limiting using Redis

### Database Per Service Pattern

Each service owns its database, ensuring:
- Independent schema evolution
- Service isolation
- Independent scaling of data layer
- No cross-service database queries

### Future Improvements: Event-Driven Architecture

Extend with Apache Kafka for:

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Auth Service   │────▶│      KAFKA      │────▶│  Task Service   │
│                 │     │                 │     │                 │
│ Events:         │     │ Topics:         │     │ Consumes:       │
│ - UserCreated   │     │ - user-events   │     │ - UserCreated   │
│ - UserUpdated   │     │ - task-events   │     │ - UserDeleted   │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

Benefits:
- Asynchronous communication
- Event sourcing
- Better fault tolerance
- Eventual consistency

## Security

### JWT Authentication
- Tokens signed with HS256 algorithm
- 15-minute token expiry
- Tokens validated on each request to Task Service

### Password Security
- Passwords hashed using bcrypt (cost factor: 12)
- No plain-text password storage
- Password complexity validation

### Environment Security
- All secrets via environment variables
- No hardcoded credentials
- Separate credentials per service

### Input Validation
- All inputs validated before processing
- SQL injection prevention via GORM
- XSS prevention through proper encoding

### CORS Configuration
- Configurable allowed origins
- Proper preflight handling
- Credentials support

## Monitoring & Health Checks

Each service exposes a `/health` endpoint:

```bash
# Check API Gateway
curl http://localhost:8080/health

# Check Auth Service
curl http://localhost:8081/health

# Check Task Service
curl http://localhost:8082/health
```

## Swagger Documentation

Access Swagger UI for each service:
- Auth Service: http://localhost:8081/swagger/index.html
- Task Service: http://localhost:8082/swagger/index.html

## Testing

```bash
# Run tests for Auth Service
cd auth-service
go test ./...

# Run tests for Task Service
cd task-service
go test ./...

# Run tests with coverage
go test -cover ./...
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
