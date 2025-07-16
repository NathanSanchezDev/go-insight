# Go Insight

This project consists of a Go backend and a Next.js frontend. Use Docker Compose to run the whole stack.

## Usage

1. Copy `.env` and adjust values as needed.
2. Run `docker compose up -d --build` from the project root.

Services:
- **backend** - Go API built from `core`
- **frontend** - Next.js UI built from `ui`
- **postgres** - Database for storing logs and metrics

See `docs/README.md` for full documentation.

