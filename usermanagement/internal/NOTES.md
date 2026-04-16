```
user-management/
├── cmd/
│   └── server/
│       └── main.go                 # Starts gRPC & HTTP servers
│
├── config/                          # Config files & loaders
│   └── config.go
│
├── bootstrap/                        # App initialization
│   └── bootstrap.go
│
├── proto/                            # gRPC protobuf definitions
│   └── user.proto
│
├── internal/
│   ├── domain/                       # Pure business logic
│   │   ├── entity/
│   │   │   ├── user.go               # User entity
│   │   │   ├── role.go               # Role entity
│   │   │   └── session.go            # Refresh token entity
│   │   │
│   │   ├── service/                  # Domain services (business rules)
│   │   │   └── pbac_service.go       # Policy evaluation
│   │   │
│   │   └── errors/                   # Domain-specific errors
│   │       └── domain_errors.go
│   │
│   ├── application/                  # Use cases / application services
│   │   ├── auth_service.go           # Registration, login, OTP
│   │   ├── session_service.go        # JWT + refresh token management
│   │   └── password_service.go       # Password reset/change use cases
│   │
│   ├── repository/                   # Persistence abstractions
│   │   ├── user_repository.go        # Interface + DB implementation
│   │   ├── session_repository.go
│   │   └── role_repository.go
│   │
│   ├── interfaces/                   # Adapters / ports
│   │   ├── grpc/                     # gRPC handlers
│   │   │   └── user_handler.go
│   │   ├── http/                     # HTTP handlers for gateway
│   │   │   └── user_handler.go
│   │   ├── dto/                      # Request / response structures
│   │   │   ├── auth_dto.go
│   │   │   ├── session_dto.go
│   │   │   └── password_dto.go
│   │   └── mapper/                   # DTO ↔ Entity conversions
│   │       └── user_mapper.go
│   │
│   ├── middleware/                   # HTTP / gRPC middlewares
│   │   └── auth_middleware.go
│   │
│   ├── utils/                         # Generic helpers
│   │   ├── jwt_utils.go
│   │   ├── otp_utils.go
│   │   └── hash_utils.go
│   │
│   └── apperrors/                     # Standardized application errors
│       └── app_errors.go
│
├── database/
│   ├── migrations/
│   ├── config/
│   └── seeding/
│
└── go.mod
```