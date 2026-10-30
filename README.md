---
name: go-hexagonal-architecture
description: Flat hexagonal architecture for Go services. Use when creating a new Go service, when adding features, adapters or configuration to an existing one, or when deciding which package something belongs in.
---

# go-hexagonal-template

An architecture template for Go services: hexagonal, flat, and limited to the standard library plus whatever drivers are needed.

The login inside is filler. It exists so the structure can be seen working. What gets reused is the architecture: this document defines the rules, and the code is an example of them.

## The idea

Business logic sits at the center and knows nothing about the outside world. HTTP, MySQL, JWT and bcrypt are details plugged in from the outside through interfaces. Replacing one of them means writing a new adapter and changing one line in `main`.

```
Client ──▶ adapters/in ──▶ ports/in ◁── service ──▶ ports/out ◁── adapters/out ──▶ MySQL, JWT...
                                           │
                                           ▼
                                         domain

──▶ calls    ◁── implements
```

`cmd/api/main.go` stays out of the picture: it is the only place that knows every piece and connects them.

## Structure

```
cmd/api/main.go                 # wiring + process lifecycle
internal/
├── config/                     # Config + Load() from environment variables (imported by main only)
├── domain/                     # entities, rules and domain errors
├── ports/                      # in.go (what it offers) / out.go (what it needs)
├── service/                    # use cases
└── adapters/
    ├── in/httpapi/             # HTTP handlers, DTOs, middleware
    └── out/
        ├── mysql/              # Open (pool with a ping at startup) + repositories
        ├── password/           # bcrypt
        └── token/              # JWT HS256
migrations/                     # schema SQL
```

It is flat on purpose: one folder per layer, with no sublayers and no package per use case. Depth is added only when the project asks for it.

## Dependency rule

This is the rule that holds everything else up: dependencies always point toward the center.

| Package | May import | Never imports |
|---|---|---|
| `domain` | stdlib | anything from the project or third parties |
| `ports` | `domain` | `service`, `adapters`, third parties |
| `service` | `domain`, `ports` | `adapters`, `config`, drivers, `net/http` |
| `adapters/in/*` | `domain`, `ports`, libraries for its technology | `service`, other adapters, `config` |
| `adapters/out/*` | `domain`, `ports`, libraries for its technology | `service`, other adapters, `config` |
| `config` | stdlib | anything from the project |
| `cmd/api` | everything | it is the only package that imports `config` and the adapters together |

If a change needs to break this table, the code belongs in another layer.

## What goes in each layer

**`domain`**: entities, their invariants and business errors. An entity is created and modified through functions that validate, so a domain value is always valid. It does not know what a database or JSON is.

**`ports`**: interfaces only, plus the input types they need.
- `in.go`: what the application **offers**. Implemented by `service`, called by inbound adapters.
- `out.go`: what the application **needs**. Called by `service`, implemented by outbound adapters.

**`service`**: use cases. Orchestrates `domain` and the outbound ports. It receives interfaces, never concrete adapter types.

**`adapters/in`**: translates a protocol into calls to the inbound ports: parse the request, call the port, serialize the response. No business logic. DTOs live here, never in `domain`.

**`adapters/out`**: implements an outbound port with a concrete technology. No business logic. Its own errors are converted into domain errors here.

**`config`**: reads environment variables, validates them and returns a struct. It is the only place that calls `os.Getenv`.

**`cmd/api/main.go`**: builds the pieces, connects them and manages the lifecycle. No logic.

## Rules

**Dependency injection**
- Everything arrives through constructors (`NewX(...)`). Adapters receive concrete values (a `dsn`, a `secret`, a `*sql.DB`), never the `Config` struct.
- No mutable global state in any package. The sentinel errors in `domain` are the only package-level variables.

**Interfaces**
- They are defined in `ports`, not next to their implementation.
- Each implementation asserts it at compile time: `var _ ports.UserRepository = (*UserRepository)(nil)`.

**Errors**
- `domain` defines the business errors.
- Outbound adapters translate their errors into those (for example, a MySQL duplicate becomes `ErrEmailTaken`).
- Inbound adapters translate those errors into their protocol in a single place (`writeError`). Anything unrecognized is logged and returned as a generic error, leaking no details to the client.

**Eager connections**
- External connections are opened and verified in `main` at startup. If one fails, the process does not start.
- No lazy loading and no `sync.Once`: they hide dependencies and delay the failure until the first request.

**`main` without a `run()` function**

Everything goes in `main()`, in this order:

1. Logger
2. `config.Load()`
3. `signal.NotifyContext`
4. Open external connections
5. Outbound adapters
6. Services
7. Inbound adapter
8. Server in a goroutine with an error channel
9. `select` between signal and error
10. `Shutdown` with a timeout
11. Close connections

`os.Exit` and `log.Fatal` do not run `defer`. Before connections are open, logging and calling `os.Exit(1)` is enough. After that, every error exit logs, closes the connections explicitly and only then calls `os.Exit(1)`.

**Tests**
- Services are tested with in-memory fakes of the ports, with no database or external libraries. If a use case cannot be tested that way, something has crossed a boundary.

## How to add things

**A new feature**, from the inside out:
1. `domain`: entity, rules and errors.
2. `ports`: the inbound port and any missing outbound ones.
3. `service`: the use case, with tests using fakes.
4. `adapters/out`: the implementation of the outbound ports.
5. `adapters/in`: the handler, its DTOs and the translation of any new errors.
6. `main`: connect the new pieces.

**A technology swap** (for example, MySQL to Postgres): a new adapter implementing the same port and a constructor change in `main`. If `service` or `domain` has to change, the port was leaking the technology.

**A new inbound adapter** (gRPC, CLI, queue consumer): a package in `adapters/in` that uses the same inbound ports.

**A configuration variable**: a field in `Config`, reading it in `Load()`, and passing the concrete value through a constructor from `main`.

## Using it as a template

1. Create the repo with **Use this template** on GitHub.
2. Change the module name:

   ```bash
   OLD=github.com/JoelTeoGom/login-service
   NEW=github.com/your-user/your-service
   go mod edit -module "$NEW"
   find . -name '*.go' -exec perl -pi -e "s#$OLD#$NEW#g" {} +
   ```

3. Change the `"login-service"` issuer in `cmd/api/main.go`.
4. Replace the example domain following **A new feature**, and delete the adapters you do not use.
5. Check that everything still holds:

   ```bash
   gofmt -l . && go vet ./... && go build ./... && go test ./...
   ```

To run the example:

```bash
docker compose up -d
export MYSQL_DSN='login:login@tcp(localhost:3306)/login?parseTime=true'
export JWT_SECRET='change-me'
go run ./cmd/api
```

## Using this README as a skill

The block at the top is Claude Code skill frontmatter, so this file works as is:

```bash
mkdir -p ~/.claude/skills/go-hexagonal-architecture
cp README.md ~/.claude/skills/go-hexagonal-architecture/SKILL.md
```
