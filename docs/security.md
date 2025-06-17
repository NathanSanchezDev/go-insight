# Security Guide

Go-Insight implements enterprise-grade security features designed for production environments. This guide covers authentication, rate limiting, and security best practices.

## Overview

Go-Insight's security model follows the principle of **defense in depth**:

1. **Authentication** - API key validation for endpoint access
2. **Rate Limiting** - Per-IP request throttling to prevent abuse
3. **Endpoint Protection** - Public monitoring vs. protected data endpoints
4. **Input Validation** - JSON schema checks and request size limits

## Authentication

### Supported Methods

Go-Insight supports multiple authentication methods for flexibility:

#### 1. X-API-Key Header (Recommended)
```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/logs
```

#### 2. Authorization Bearer Token
```bash
curl -H "Authorization: Bearer your-api-key" http://localhost:8080/logs
```

#### 3. Authorization ApiKey
```bash
curl -H "Authorization: ApiKey your-api-key" http://localhost:8080/logs
```

#### 4. Query Parameter (Development Only)
```bash
curl "http://localhost:8080/logs?api_key=your-api-key"
```

**⚠️ Note**: Query parameter authentication should only be used in development. Use header-based authentication in production.

### Configuration

Set your API key in the environment:

```bash
# .env file
API_KEY=your-secure-api-key-here
```

**Best Practices for API Keys:**
- Use a strong, randomly generated key (minimum 32 characters)
- Rotate keys regularly (recommended: every 90 days)
- Store keys securely (environment variables, not in code)
- Use different keys for different environments

### Error Responses

Invalid or missing authentication returns helpful error messages:

```json
{
  "error": "API key required",
  "hint": "Provide API key in Authorization header, X-API-Key header, or api_key query parameter"
}
```

```json
{
  "error": "Invalid API key"
}
```

## Rate Limiting

### Implementation

Go-Insight implements **per-IP rate limiting** to prevent API abuse:

- **Default Limit**: 60 requests per minute per IP address
- **Tracking Method**: Token bucket algorithm with sliding window
- **Scope**: Rate limits apply per client IP address
- **Headers**: Real-time rate limit information in response headers

### Configuration

Configure rate limiting via environment variables:

```bash
# .env file
RATE_LIMIT_REQUESTS=60    # Requests per window (default: 60)
RATE_LIMIT_WINDOW=1       # Window in minutes (default: 1)
```

### Rate Limit Headers

All API responses include rate limiting information:

```http
X-RateLimit-Limit: 60           # Maximum requests per window
X-RateLimit-Remaining: 45       # Requests remaining in current window
```

When rate limited, additional headers are included:

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
Retry-After: 60                 # Seconds until reset
```

### Rate Limit Response

```json
{
  "error": "Rate limit exceeded",
  "limit": 60,
  "window": "1 minute"
}
```

### IP Address Detection

The rate limiter intelligently detects client IP addresses:

1. **X-Forwarded-For** header (for load balancers/proxies)
2. **X-Real-IP** header (for nginx/other proxies)  
3. **Remote address** (direct connections)

This ensures accurate rate limiting behind reverse proxies and load balancers.

### Per-IP Isolation

Each IP address maintains its own rate limit bucket:

- **Separate Tracking**: Different IPs don't affect each other's limits
- **Memory Efficient**: Automatic cleanup of expired IP entries
- **Thread Safe**: Concurrent requests from same IP handled correctly

## Endpoint Protection

### Public Endpoints

These endpoints are accessible without authentication:

- `GET /health` - System health check (for monitoring systems)

### Protected Endpoints

All data endpoints require valid API key authentication:

- `GET /logs` - Log retrieval
- `POST /logs` - Log ingestion
- `GET /metrics` - Metrics retrieval  
- `POST /metrics` - Metrics ingestion
- `GET /traces` - Trace retrieval
- `POST /traces` - Trace creation
- All other trace and span endpoints

### Middleware Order

Security middleware is applied in the following order:

1. **Rate Limiting** - Applied first to prevent resource exhaustion
2. **Authentication** - API key validation for protected endpoints
3. **Request Logging** - Audit trail of all requests

## Security Headers

Go-Insight includes security-relevant headers in responses:

```http
# Rate limiting information
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45

# Content type specification
Content-Type: application/json

# CORS headers (when configured)
Access-Control-Allow-Origin: *
```

## Production Security Checklist

### Environment Security
- [ ] API keys stored in environment variables, not code
- [ ] Database credentials secured and rotated regularly
- [ ] TLS/HTTPS enabled for all communications
- [ ] Firewall rules restrict database access

### Application Security  
- [ ] Strong API keys (32+ characters, randomly generated)
- [ ] Rate limiting configured appropriately for your traffic
- [ ] Monitoring and alerting for authentication failures
- [ ] Regular security updates and dependency scanning

### Network Security
- [ ] Go-Insight deployed behind reverse proxy (nginx, Apache)
- [ ] Load balancer configured with proper IP forwarding
- [ ] Database not directly exposed to internet
- [ ] VPC/network segmentation implemented

### Monitoring Security
- [ ] Authentication failures logged and monitored
- [ ] Rate limit violations tracked and alerted
- [ ] Unusual traffic patterns detected
- [ ] Security incidents have response procedures

## Common Security Scenarios

### Behind Load Balancer/Proxy

When deploying behind a load balancer, ensure IP forwarding is configured:

```nginx
# nginx configuration
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Real-IP $remote_addr;
```

### Multiple API Keys

For different applications or environments, use environment-specific API keys:

```bash
# Development
API_KEY=dev-key-12345

# Staging  
API_KEY=staging-key-67890

# Production
API_KEY=prod-key-abcdef
```

### Rate Limit Bypass

For trusted internal services, you can implement IP whitelisting by modifying the rate limit middleware to skip certain IP ranges.

## Security Monitoring

### Logging

Security events are logged for monitoring:

```
🔒 Authentication failed: No API key provided from GET /logs
🔒 Authentication failed: Invalid API key from POST /metrics  
🚦 Rate limit exceeded for IP 192.168.1.100 on GET /logs
✅ Authenticated request: GET /logs
```

### Metrics to Monitor

- Authentication failure rate
- Rate limit violations per IP
- Unusual traffic patterns  
- Geographic distribution of requests
- API endpoint usage patterns

## Future Security Enhancements

Planned security improvements include:

- **Security Headers**: CSRF protection and security-focused HTTP headers
- **Request/Response Encryption**: End-to-end encryption for sensitive data
- **Audit Logging**: Comprehensive security event tracking

## Support

For security-related questions or to report security vulnerabilities:

- Create an issue on GitHub for general security questions
- For security vulnerabilities, please email nathansanchezdev@outlook.com

---

**Last Updated**: May 2025  
**Version**: 1.0.0
