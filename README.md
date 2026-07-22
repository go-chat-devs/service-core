# Service Core

Core service of the Chat Application — handles user storage, database connectivity, and structured logging.

## Tech Stack

- **Go** 1.25.4
- **PostgreSQL** (via [pgx/v5](https://github.com/jackc/pgx) connection pool)
- **Structured logging** ([log/slog](https://pkg.go.dev/log/slog) + [fatih/color](https://github.com/fatih/color))

## Requirements

- Go ≥ 1.25
- PostgreSQL

## Project Structure

```
service-core/
├── cmd/
│   └── service-core/
│       └── main.go            # Application entrypoint
├── internal/
│   ├── logger/
│   │   ├── logger.go          # Default slog logger setup
│   │   └── handler.go         # Pretty-printed, color-coded slog handler
│   └── storage/
│       ├── storage.go         # Root storage — PostgreSQL pool init
│       └── users/
│           └── storage.go     # User-specific storage layer
├── go.mod
├── go.sum
└── .gitignore
```

## Environment Variables

All variables live in `.env` at the project root. Create it before running:

```bash
cp .env.example .env
```

| Variable | Example                                        | Description                     |
| -------- | ---------------------------------------------- | ------------------------------- |
| `DB_URL`  | `postgres://user:pass@localhost:5432/chat?sslmode=disable` | PostgreSQL connection string    |

## Running

```bash
# Clone the repository
git clone https://github.com/go-chat-devs/service-core.git
cd service-core

# Install dependencies
go mod download

# Run the service
go run ./cmd/service-core
```

### Build

```bash
go build -o bin/service-core ./cmd/service-core
./bin/service-core
```

# Authors

- @Pachandre [github](https://github.com/Pachandre)
