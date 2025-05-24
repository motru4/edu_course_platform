# Auth Service

Authentication and authorization service written in Golang with JWT, OAuth 2.0, and 2FA support.

## 🚀 Features

- User registration with email validation and password verification
- JWT authentication (access + refresh tokens)
- User roles (student/author/admin)
- !TODO OAuth 2.0 (Google) 
- 2FA via email
- gRPC service for access verification
- Rate limiting
- Enhanced password security with bcrypt + pepper

## 📋 Requirements

- Go 1.21+
- PostgreSQL 14+
- Make

## 🛠 Installation

1. Clone the repository:
```bash
git clone https://github.com/motru4/auth-service.git
cd auth-service
```

2. Install dependencies:
```bash
make deps
```

3. Create .env file based on .env.example:
```bash
cp .env.example .env
```

4. Configure environment variables in .env file:
```env
# Required security settings
PASSWORD_PEPPER=your-secure-pepper    # Used for password hashing
JWT_SECRET=your-jwt-secret           # Used for token signing

# Other settings
DB_URL=postgresql://localhost:5432/auth_service
```

5. Create database:
```bash
createdb auth_service
```

6. Apply migrations:
```bash
make migrate-up
```

## 🚀 Running

```bash
# Run service
make run
```

## 📝 API Endpoints

### HTTP API (8080)

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - Authentication
- `POST /api/v1/auth/verify-email` - Email verification for registration
- `POST /api/v1/auth/verify-login` - 2FA verification for login
- `POST /api/v1/auth/refresh` - Refresh token pair
- `POST /api/v1/auth/reset-password/request` - Request password reset
- `POST /api/v1/auth/reset-password/confirm` - Confirm password reset
- `GET /api/v1/auth/oauth/google` - OAuth 2.0 via Google
- `GET /api/v1/auth/swagger/*` - API documentation

### gRPC API (9090)

- `CheckAccess` - Token access verification

## 🔒 Security

- Password security:
  - Bcrypt hashing with configurable work factor
  - Additional pepper for enhanced security
  - Separate pepper storage from database
- JWT tokens:
  - Signed using HS256
  - Automatic invalidation on password change
  - Configurable expiration times
- Rate limiting: 5 requests per minute
- Prepared statements for SQL injection protection
- Automatic cleanup of expired refresh tokens
- Email verification for registration

## 📦 Project Structure

```
.
├── cmd/                   # Application entry points
├── docs/                  # Project documentation
├── internal/
│   ├── app/              # Application initialization
│   ├── config/           # Configuration
│   ├── database/         # Database operations
│   ├── models/           # Data models
│   ├── repositories/     # Repositories
│   ├── services/         # Business logic
│   ├── security/         # Security utilities (password hashing, etc.)
│   ├── transport/        # API (HTTP + gRPC)
│   └── utils/            # Helper functions
└── migrations/           # SQL migrations
```

## 📄 License

MIT 
