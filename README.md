# go-ez-auth

## Overview
Unified, modular, and extensible authentication library for Go web applications, inspired by Passport.js and Django Auth.

## Why go-ez-auth?
go-ez-auth provides a unified, modular, and extensible authentication framework for Go web applications. It reduces boilerplate, enforces security best practices, and works across multiple web frameworks.

## Key Features
- Modular authentication strategies: JWT, Session, API Key, OAuth2/OIDC, Local (username/password)
- Framework middleware adapters: net/http, Gin, Echo (extendable to others like Fiber)
- Security features: CSRF protection, secure session cookies, session fixation prevention, bcrypt password hashing
- Strategy registry: central registration and lookup of auth strategies
- Comprehensive unit tests for all components
- Quickstart examples to get up and running in minutes

## Implementation Details
- Core interfaces in `core` for Strategy, User, UserStore, and standardized errors
- Strategy implementations in `strategies/...` leveraging Gorilla sessions, golang-jwt, oauth2, bcrypt
- Middleware adapters in `middleware/...` for seamless integration with net/http, Gin, and Echo
- Examples and tests organized per strategy/middleware directory

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

## Phase 4: Security Enhancements & Features (Completed)

- **CSRF Protection**: `middleware.CSRFMiddleware` using `github.com/gorilla/csrf` for net/http. Tests in `middleware/csrf_test.go`.
- **Session Security**: HTTP-only, Secure, and SameSite flags set by default in `strategies/session` Setup; session fixation via `Login` method. Tests in `strategies/session`.
- **OAuth2 / OIDC Strategy**: Implemented in `strategies/oauth2`; uses `golang.org/x/oauth2`, userinfo fetch, and customizable extractor. Tests in `strategies/oauth2/oauth2_test.go`.
- **Local Strategy (username/password)**: Implemented in `strategies/local`; Basic Auth with bcrypt via `golang.org/x/crypto/bcrypt`. Tests in `strategies/local/local_test.go`.

**Test Phase 4**
```bash
go test ./middleware -v
go test ./strategies/session -v
go test ./strategies/oauth2 -v
go test ./strategies/local -v
```

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

### Quickstart (net/http)
```go
package main

import (
   "net/http"
   "github.com/gorilla/sessions"
   "go-ez-auth/core"
   "go-ez-auth/strategies/jwt"
   "go-ez-auth/middleware"
)

func main() {
   // Register JWT strategy
   store := sessions.NewCookieStore([]byte("secret"))
   strat := jwt.New(jwt.Config{Secret: []byte("mysecret"), UserStore: core.InMemoryUserStore{"u1"}})
   core.RegisterStrategy(strat)

   // Protected handler
   handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       user := r.Context().Value(core.ContextUserKey).(core.User)
       w.Write([]byte("Hello, " + user.GetID()))
   })

   http.Handle("/", middleware.Middleware("jwt")(handler))
   http.ListenAndServe(":8080", nil)
}
```

### Quickstart (Gin)
```go
package main

import (
   "github.com/gin-gonic/gin"
   "github.com/gorilla/sessions"
   "go-ez-auth/core"
   "go-ez-auth/strategies/jwt"
   "go-ez-auth/middleware"
)

func main() {
   r := gin.Default()
   store := sessions.NewCookieStore([]byte("secret"))
   strat := jwt.New(jwt.Config{Secret: []byte("mysecret"), UserStore: core.InMemoryUserStore{"u1"}})
   core.RegisterStrategy(strat)
   r.Use(middleware.GinMiddleware("jwt"))
   r.GET("/", func(c *gin.Context) {
       user := c.MustGet(core.ContextUserKey).(core.User)
       c.String(200, "Hello, %s", user.GetID())
   })
   r.Run(":8080")
}
```

### Quickstart (Echo)
```go
package main

import (
   "github.com/labstack/echo/v4"
   "github.com/gorilla/sessions"
   "go-ez-auth/core"
   "go-ez-auth/strategies/jwt"
   "go-ez-auth/middleware"
)

func main() {
   e := echo.New()
   store := sessions.NewCookieStore([]byte("secret"))
   strat := jwt.New(jwt.Config{Secret: []byte("mysecret"), UserStore: core.InMemoryUserStore{"u1"}})
   core.RegisterStrategy(strat)
   e.Use(middleware.EchoMiddleware("jwt"))
   e.GET("/", func(c echo.Context) error {
       user := c.Get(core.ContextUserKey).(core.User)
       return c.String(200, "Hello, %s", user.GetID())
   })
   e.Start(":8080")
}
```

Run all tests:
```bash
go test ./... -cover
```

## Contributing

Contributions welcome. Please open issues or PRs against this repo.

## License

MIT
