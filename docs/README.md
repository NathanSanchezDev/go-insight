# Go-Insight

<p align="center">
  <img src="https://github.com/user-attachments/assets/b6f57e12-9170-490a-9d2c-8b0359b85e47" alt="Go-Insight Logo" width="200" height="200">
</p>

<p align="center">
  <a href="#about">About</a> •
  <a href="#features">Features</a> •
  <a href="#getting-started">Getting Started</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#api-endpoints">API Endpoints</a> •
  <a href="#roadmap">Roadmap</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

## About
Go-Insight is a modern observability platform designed for distributed applications, built with Go and PostgreSQL. It collects, stores, and analyzes logs, metrics, and distributed traces in one unified system, giving developers comprehensive visibility into their application performance and behavior.

Unlike heavyweight enterprise observability solutions, Go-Insight focuses on being:

- **Lightweight**: Minimal resource requirements, optimal for small to medium deployments
- **Self-hosted**: Full control over your observability data without vendor lock-in
- **Developer-friendly**: Simple API and SDK integration for any application
- **Extensible**: Modular architecture that can grow with your needs

## Features

### Current Features (Phase 1)

- **✅ Comprehensive Logging**
  - Structured log ingestion with support for log levels, service names, and metadata
  - Advanced filtering and querying capabilities
  - Correlation with traces via trace IDs and span IDs

- **✅ Performance Metrics**
  - HTTP endpoint metrics (response times, status codes, etc.)
  - Custom metadata fields for framework, language, and version
  - Flexible querying with support for multiple filters

- **✅ Distributed Tracing**
  - Full request lifecycle tracking across services
  - Parent-child relationship modeling with spans
  - Timing and duration measurements

- **✅ Developer Experience**
  - Easy-to-use RESTful API
  - Containerized setup with Docker Compose
  - Comprehensive API documentation

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (for containerized setup)

### Environment Setup

Create a `.env` file in the project root with the following variables:

```
DB_USER=postgres
DB_PASS=yourpassword
DB_NAME=goinsight
DB_HOST=localhost
DB_PORT=5432
PORT=8080
```

### Running with Docker Compose

The quickest way to get started is with Docker Compose:

```bash
# Clone the repository
git clone https://github.com/NathanSanchezDev/go-insight.git
cd go-insight

# Start the database and application
docker-compose up -d
```

### Manual Setup

```bash
# Clone the repository
git clone https://github.com/NathanSanchezDev/go-insight.git
cd go-insight

# Set up the database
./scripts/setup_db.sh

# Build and run
go build -o go-insight ./cmd/main.go
./go-insight
```

### Verifying Installation

Once running, you can verify the setup by accessing the health endpoint:

```bash
curl http://localhost:8080/health
# Should return "OK"
```

## Architecture

Go-Insight follows a modular architecture with several key components:

- **API Layer**: RESTful endpoints for data ingestion and retrieval
- **Storage Layer**: PostgreSQL for durable storage of all observability data
- **Service Layer**: Business logic for processing and correlating data
- **SDK Layer** (coming soon): Client libraries for popular languages

### Database Schema

The system uses three main tables:

- **logs**: Stores structured log entries with metadata
- **metrics**: Captures performance metrics from API endpoints
- **traces** and **spans**: Tracks distributed request execution across services

For detailed schema information, see the migration files in `internal/db/migrations/`.

## API Endpoints

### Logs API

- **GET /logs**: Fetch logs with filtering options
  - Query params: `service`, `level`, `message`, `start_time`, `end_time`, `limit`, `offset`
- **POST /logs**: Ingest new log entries

### Metrics API

- **GET /metrics**: Fetch metrics with filtering options
  - Query params: `service`, `path`, `method`, `min_status`, `max_status`, `limit`, `offset`
- **POST /metrics**: Ingest new metrics

### Tracing API

- **GET /traces**: Fetch traces with optional filtering
- **POST /traces**: Create a new trace
- **POST /traces/{traceId}/end**: Mark a trace as completed
- **GET /traces/{traceId}/spans**: Get all spans for a specific trace
- **POST /spans**: Create a new span
- **POST /spans/{spanId}/end**: Mark a span as completed

## Roadmap

Go-Insight is under active development. See our detailed [roadmap](docs/roadmap.md) for upcoming features, including:

### Phase 2: Enhanced Core Features (1-2 months)
- API improvements for high-volume data
- Distributed tracing with context propagation
- Data retention and compression

### Phase 3: User Interface (2-3 months)
- Web UI for visualizing metrics and logs
- Customizable dashboards
- Alerting system

### Phase 4: Advanced Features (3-4 months)
- Authentication and access control
- Multi-language client SDKs
- Advanced analytics and anomaly detection

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/NathanSanchezDev">Nathan Sanchez</a>
</p>
