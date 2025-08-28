# User Service

A microservice for user authentication and management built with Go, gRPC, and PostgreSQL.

## 🚀 Features

- **User Authentication**: JWT-based authentication with access and refresh tokens
- **User Management**: CRUD operations for users
- **Role Management**: Role-based access control (RBAC)
- **Session Management**: Secure session handling with Redis
- **Email Integration**: gRPC-based email service integration
- **Queue System**: Asynchronous task processing with Redis
- **Service Discovery**: Automatic service registration and discovery
- **Database Migrations**: Automated database schema management

## 🏗️ Architecture

The service follows a clean architecture pattern with the following layers:

```
├── cmd/                    # Application entry point
├── bootstrap/             # Application initialization
├── domain/               # Business logic and entities
│   ├── entity/          # Domain entities
│   ├── repository/      # Repository interfaces
│   └── usecase/         # Business use cases
├── infrastructure/       # External dependencies
│   ├── grpc_client/     # gRPC client implementations
│   ├── grpc_service/    # gRPC server implementation
│   └── repo/            # Repository implementations
├── constants/           # Application constants
└── migrations/          # Database migrations
```

## 📋 Prerequisites

- Go 1.24.6 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Docker (optional, for development)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd user-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```bash
   # Create database
   make dev-create-db
   
   # Run migrations
   make migrate-dev-up
   ```

4. **Configure environment**
   Copy the development configuration:
   ```bash
   cp dev.config.yaml config.yaml
   ```
   
   Update the configuration with your database and Redis credentials.

## 🚀 Running the Service

### Development Mode
```bash
# Run the service
make run

# Or directly with go
go run cmd/main.go
```

### Production Build
```bash
# Build the binary
make build-grpc

# Run the binary
./bin/grpc-server
```

## 📊 Database Management

### Migrations
```bash
# Apply migrations
make migrate-dev-up

# Rollback migrations
make migrate-dev-down

# Reset database (drop all tables and reapply migrations)
make migrate-dev-reset

# Create new migration
make migrate-dev-create name=migration_name
```

### Database Operations
```bash
# Drop database
make dev-drop-db

# Create database
make dev-create-db
```

## 🔧 Configuration

The service uses a YAML configuration file. Key configuration options:

```yaml
MODE_ENV: "development"                    # Environment mode
URL_DB: "postgres://..."                   # Database connection string
NAME_SERVICE: "AuthService"                # Service name for discovery
PORT_GRPC: 50050                          # gRPC server port
HOST_GRPC: "localhost"                     # gRPC server host
INTERVAL_CHECK: "10s"                     # Health check interval

# Redis configuration for cache and queue
DB_CACHE:
  Addr: "localhost:6379"
  DB: 1
  Password: ""
  
QUEUE:
  Addr: "localhost:6379"
  DB: 0
  Password: ""

# JWT secrets
JWT_SECRET:
  Access: "your-access-jwt-secret"
  Refresh: "your-refresh-jwt-secret"
  Verify: "your-verify-jwt-secret"
  Forgot: "your-forgot-jwt-secret"

# External services
MAIL_SERVICE_ADDR: "localhost:50052"
```

## 🔌 API Endpoints

The service exposes gRPC endpoints for:

### User Management
- `CreateUser` - Create a new user
- `GetUser` - Retrieve user information
- `UpdateUser` - Update user details
- `DeleteUser` - Delete a user
- `GetAllUsers` - List all users

### Authentication
- `Login` - User login with credentials
- `Register` - User registration
- `RefreshToken` - Refresh access token
- `Logout` - User logout
- `VerifyEmail` - Email verification
- `ForgotPassword` - Password reset request
- `ResetPassword` - Password reset

### Role Management
- `CreateRole` - Create a new role
- `GetRole` - Get role by ID
- `GetRoleByName` - Get role by name
- `UpdateRole` - Update role details
- `DeleteRole` - Delete a role
- `GetAllRoles` - List all roles
- `CheckRole` - Check user role permissions

### Session Management
- `GetSessions` - Get user sessions
- `DeleteSession` - Delete specific session
- `DeleteExpiredSessions` - Clean up expired sessions

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./domain/usecase/user
```

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MODE_ENV` | Environment mode | `development` |
| `URL_DB` | Database connection string | - |
| `NAME_SERVICE` | Service name | `AuthService` |
| `PORT_GRPC` | gRPC server port | `50050` |
| `HOST_GRPC` | gRPC server host | `localhost` |
| `INTERVAL_CHECK` | Health check interval | `10s` |
| `SECRET_OTP` | OTP secret key | - |
| `JWT_SECRET_ACCESS` | JWT access token secret | - |
| `JWT_SECRET_REFRESH` | JWT refresh token secret | - |
| `FRONTEND_URL` | Frontend URL | `http://localhost:3000` |

## 🔒 Security

- JWT tokens with configurable expiration
- Password hashing with bcrypt
- Role-based access control (RBAC)
- Session management with Redis
- Input validation and sanitization
- Secure communication over gRPC

## 📦 Dependencies

### Core Dependencies
- `github.com/anhvanhoa/service-core` - Core service utilities
- `github.com/go-pg/pg/v10` - PostgreSQL ORM
- `go.uber.org/zap` - Structured logging
- `google.golang.org/grpc` - gRPC framework
- `github.com/redis/go-redis/v9` - Redis client

### Development Dependencies
- `github.com/onsi/ginkgo` - Testing framework
- `github.com/onsi/gomega` - Testing matchers

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation

## 🔄 Changelog

### Version 1.0.0
- Initial release
- User authentication and management
- Role-based access control
- Session management
- Email integration
- Queue system integration
