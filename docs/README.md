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

## Quick Start

### 1. Environment Setup
```bash
# Clone and setup
git clone https://github.com/NathanSanchezDev/go-insight.git
cd go-insight

# Configure environment
# Copy the sample file from the repository root
cp .env.example .env
# Edit .env with your database credentials and API key
```

### 2. Database Setup
```bash
# Start PostgreSQL (Docker)
docker-compose up -d

# Run migrations  
./scripts/setup_db.sh
```

### 3. Start Server
```bash
# Run application
go run cmd/main.go

# Verify health
curl http://localhost:8080/health
# Returns: OK
```

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
  -d '{"service_name": "api-service", "path": "/login", "method": "POST", "status_code": 200, "duration_ms": 45.7}'
```

### Query Data
```bash
# Filter logs by service
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/logs?service=api-service&level=ERROR&limit=10"

# Get performance metrics
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/metrics?service=api-service&min_status=400"
```

## Performance Benchmarks

| Operation | Response Time | Throughput |
|-----------|---------------|------------|
| Service-based log queries | ~5ms | 1000+ req/sec |
| Metrics with complex filters | ~5ms | 800+ req/sec |  
| Trace lookups | ~5ms | 1200+ req/sec |
| Concurrent connections | <1ms overhead | 100+ simultaneous |

## Documentation

- **[Security Guide](security.md)** - Authentication, rate limiting, and security best practices
- **[API Reference](api.md)** - Complete endpoint documentation with examples
- **[Performance Guide](performance.md)** - Optimization strategies and benchmarking
- **[Deployment Guide](deployment.md)** - Production deployment and configuration
- **[Architecture Overview](architecture.md)** - System design and database schema

## Roadmap

**🎯 Phase 1: Foundation** ✅ *Complete*  
Core APIs, security, and performance optimization

**🚀 Phase 2: Production Features** *In Progress*
Input validation (completed), internal monitoring, enhanced logging

**📊 Phase 3: User Interface** *Planned*  
Web dashboard, visualizations, alerting system

**🔧 Phase 4: Advanced Features** *Future*  
Bulk APIs, multi-tenancy, client SDKs

[View detailed roadmap →](roadmap.md)

## Requirements

- **Go** 1.23+
- **PostgreSQL** 15+
- **Docker** (optional, for easy setup)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>⭐ Star this repo if Go-Insight helps with your observability needs!</strong><br>
  Made with ❤️ by <a href="https://github.com/NathanSanchezDev">Nathan Sanchez</a>
</p>