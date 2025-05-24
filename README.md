# Educational Platform

A microservice-based online learning platform with support for courses, lessons, and learning progress tracking.

## Architecture

The project consists of two main microservices:

### Auth Service
- User management
- Authentication and authorization
- JWT tokens
- Email verification
- Password management

### Edu Service
- Course and lesson management
- Course categories
- Learning progress tracking
- Testing system
- Content moderation
- Payment system

## Technology Stack

- **Programming Language:** Go
- **Database:** PostgreSQL
- **API:** REST API + gRPC
- **Documentation:** Swagger
- **Containerization:** Docker
- **Migrations:** Goose

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.19+
- PostgreSQL
- Make (for auth service)

### Installation and Launch

1. Clone the repository:
```bash
git clone [repository URL]
```

2. Create environment files from examples:
```bash
cp auth/.env.example auth/.env
cp edu/.env.example edu/.env
cp docker-compose.yml.example docker-compose.yml
```

3. Configure environment variables in `.env` and `docker-compose.yml` files

4. Start services using Docker Compose:
```bash
docker-compose up -d
```

### Migrations

Database migrations are located in the `/migrations` directory. Goose is used for applying migrations.

## API Documentation

- Auth Service Swagger: `http://localhost:[auth-port]/swagger/index.html`
- Edu Service Swagger: `http://localhost:[edu-port]/swagger/index.html`

## Project Structure

```
.
├── auth/                 # Authentication service
├── edu/                  # Educational service
└── migrations/           # Database migrations
```

## License

See [LICENSE](LICENSE) file 