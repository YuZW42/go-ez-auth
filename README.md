# go-ez-auth

Unified, modular, and extensible authentication library for Go web applications, inspired by Passport.js and Django Auth.

## Core Principles

- **Modularity**: Plug in or swap authentication strategies easily.
- **Extensibility**: Define custom strategies.
- **Framework Agnosticism**: Core logic is independent of web frameworks.
- **Ease of Use**: Simple middleware adapters for popular frameworks.
- **Security**: Sensible defaults and best practices.

## Project Structure

```text
/go-ez-auth
├── core        # Core interfaces, registry, errors
├── stores      # Default UserStore implementations
├── strategies  # Authentication strategy implementations
├── middleware  # Framework middleware adapters
├── examples    # Sample applications
├── internal    # Internal utilities
├── go.mod
└── README.md
```

## Phase 1: Core Architecture & Foundation (Completed)

- Defined `Strategy`, `User`, `UserStore` interfaces in `core/core.go`
- Implemented strategy registry (`RegisterStrategy`, `GetStrategy`, `ListStrategies`)
- Standard errors (`ErrUnauthorized`, `ErrInvalidCredentials`, `ErrUserNotFound`)
- Default in-memory store in `stores/inmemory.go`
- Unit tests for core and store (100% coverage)

**Test Phase 1**
```bash
go test ./core
go test ./stores
go test ./... -cover
```

## Phase 2: Core Strategies (Completed)

- **JWT Strategy** (`strategies/jwt`): token-based auth with `github.com/golang-jwt/jwt/v5`.
- **Session Strategy** (`strategies/session`): cookie-based sessions using `github.com/gorilla/sessions`.
- **API Key Strategy** (`strategies/apikey`): header & query-param lookup with `stores.APIKeyStore`.
- *(Stretch)* Basic Auth Strategy

**Test Phase 2**
```bash
go test ./strategies/jwt -v
go test ./strategies/session -v
go test ./strategies/apikey -v
```

## Phase 3: Middleware Adapters (Completed)

- Generic `AuthenticateRequest` in `middleware/nethttp.go`.
- `net/http` middleware implementation & tests.
- Gin middleware adapter & tests.
- Echo middleware adapter & tests.
- *(Optional)* Fiber adapter

**Test Phase 3**
```bash
go test ./middleware -v
```

## Phase 4: Security Enhancements & Features

- CSRF protection
- Advanced session security (fixation, flags)
- OAuth2 / OIDC strategy
- Password hashing (bcrypt) for local strategy

## Phase 5: Documentation, Tooling & Examples

- Comprehensive README and Quickstart guides
- GoDoc comments for all public APIs
- Example apps per framework in `examples/`
- *(Stretch)* Scaffolding CLI tool

## Getting Started

Install:
```bash
go get github.com/yourusername/go-ez-auth
```

Import & use core:
```go
import "go-ez-auth/core"
```

Run all tests:
```bash
go test ./... -cover
```

## Contributing

Contributions welcome. Please open issues or PRs against this repo.

## License

MIT
