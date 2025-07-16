# Go-Insight

<p align="center">
  <img src="https://github.com/user-attachments/assets/149f2e32-daad-4222-a705-df80332b1738" alt="Go-Insight Logo" width="400" height="400">
</p>

<p align="center">
  <strong>Modern observability platform for distributed applications</strong><br>
  Built with Go and PostgreSQL for production workloads
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#api-usage">API Usage</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#roadmap">Roadmap</a>
</p>

---

## Why Go-Insight?

**🚀 Production-Ready** • Enterprise-grade security, rate limiting, and performance optimization  
**⚡ Lightning Fast** • Sub-10ms query response times with strategic database indexing  
**🔒 Secure by Default** • API authentication, per-IP rate limiting, and comprehensive input validation  
**🏗️ Self-Hosted** • Full control over your observability data without vendor lock-in  
**📊 Complete Observability** • Logs, metrics, and distributed traces in one unified platform  

## Features

### ✅ Core Observability
- **Structured Logging** with filtering, search, and trace correlation
- **Performance Metrics** collection for HTTP endpoints and custom events  
- **Distributed Tracing** with parent-child span relationships

### ✅ Production Security
- **Multi-Method Authentication** (API keys, Bearer tokens)
- **Per-IP Rate Limiting** (60 req/min) with live headers
- **Public/Protected Endpoints** for monitoring and data access

### ✅ Performance Optimized
- **Sub-10ms Response Times** for filtered queries
- **Strategic Database Indexing** for 10-100x performance improvements
- **Concurrent Request Handling** with thread-safe middleware

### ✅ Container-Native
- **Docker Compose** setup with automated migrations
- **Container orchestration** ready for Kubernetes deployment
- **Environment-based configuration** for different deployment scenarios

## Quick Start

### 1. Environment Setup
```bash
# Clone and setup
git clone https://github.com/NathanSanchezDev/go-insight.git
cd go-insight

# Configure environment
cp .env.example .env
# Edit .env with your database credentials and API key
```

### 2. Run with Docker (Recommended)
```bash
# Start the entire stack
docker compose up -d --build

# Verify health
curl http://localhost:8080/health
# Returns: OK
```

### 3. Test API
```bash
# Send a test log
curl -X POST http://localhost:8080/logs \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"service_name": "test-service", "log_level": "INFO", "message": "Hello Go-Insight!"}'

# Query logs
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/logs?limit=5"
```

### Alternative: Local Development
```bash
# If you prefer to run Go locally (requires local PostgreSQL)
go run cmd/main.go
```

## Deployment

### Docker Compose (Production-Ready)
Go-Insight includes a complete Docker Compose setup with:

- **Automated migrations** - Database schema applied on startup
- **Container networking** - Services communicate via Docker network
- **Environment configuration** - Configurable via .env file
- **Data persistence** - PostgreSQL data persisted across restarts

```bash
# View container status
docker compose ps

# View application logs  
docker compose logs backend

# View database logs
docker compose logs postgres

# Stop all services
docker compose down
```

### Kubernetes Deployment
*Coming soon - Helm charts for Kubernetes deployment*

## Configuration

Go-Insight uses environment variables for configuration:

```bash
# Database settings
DB_USER=postgres
DB_PASS=your_secure_password
DB_NAME=go_insight
DB_HOST=postgres  # Use 'localhost' for local development
DB_PORT=5432

# Application settings
PORT=8080
API_KEY=your_secure_api_key

# Rate limiting
RATE_LIMIT_REQUESTS=60
RATE_LIMIT_WINDOW=1
```

**Note:** When running with Docker, use `DB_HOST=postgres`. For local development, use `DB_HOST=localhost`.

## API Usage

### Authentication
```bash
# All data endpoints require API key
curl -H "X-API-Key: your-api-key" http://localhost:8080/logs

# Rate limiting headers included in responses
# X-RateLimit-Limit: 60
# X-RateLimit-Remaining: 59
```

### Ingest Data
```bash
# Send logs
curl -X POST http://localhost:8080/logs \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"service_name": "api-service", "log_level": "INFO", "message": "User login successful"}'

# Send metrics  
curl -X POST http://localhost:8080/metrics \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"service_name": "api-service", "path": "/login", "method": "POST", "status_code": 200, "duration_ms": 45.7, "source": {"language": "go", "framework": "gin", "version": "1.9.1"}}'
```

### Query Data
```bash
# Filter logs by service
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/logs?service=api-service&level=ERROR&limit=10"

# Get performance metrics
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/metrics?service=api-service&min_status=400"

# Get distributed traces
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/traces?service=api-service&limit=10"
```

## Performance Benchmarks

| Operation | Response Time | Throughput |
|-----------|---------------|------------|
| Service-based log queries | ~5ms | 1000+ req/sec |
| Metrics with complex filters | ~5ms | 800+ req/sec |  
| Trace lookups | ~5ms | 1200+ req/sec |
| Concurrent connections | <1ms overhead | 100+ simultaneous |

## Documentation

- **[Security Guide](docs/security.md)** - Authentication, rate limiting, and security best practices
- **[API Reference](docs/api.md)** - Complete endpoint documentation with examples
- **[Performance Guide](docs/performance.md)** - Optimization strategies and benchmarking
- **[Architecture Overview](docs/architecture.md)** - System design and database schema
- **[Usage Guide](docs/usage.md)** - Integration examples and best practices

## Roadmap

**🎯 Phase 1: Foundation** ✅ *Complete*  
Core APIs, security, and performance optimization

**🚀 Phase 2: Production Features** ✅ *Complete*
Containerization, automated migrations, input validation

**📊 Phase 3: Kubernetes & Monitoring** *In Progress*  
Helm charts, Prometheus integration, Grafana dashboards

**🔧 Phase 4: Advanced Features** *Future*  
Kafka integration, multi-tenancy, client SDKs

[View detailed roadmap →](docs/roadmap.md)

## Requirements

### Docker Deployment (Recommended)
- **Docker** & **Docker Compose**
- **2GB RAM** for containers
- **1GB disk space** for data persistence

### Local Development  
- **Go** 1.23+
- **PostgreSQL** 15+
- **1GB RAM** for development

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>⭐ Star this repo if Go-Insight helps with your observability needs!</strong><br>
  Made with ❤️ by <a href="https://github.com/NathanSanchezDev">Nathan Sanchez</a>
</p>
