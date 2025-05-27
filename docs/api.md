# API Reference

Complete reference for Go-Insight's REST API. All endpoints require authentication except where noted.

## Base URL

```
http://localhost:8080
```

## Authentication

All API requests (except health check) require authentication via API key.

### Authentication Methods

| Method | Header | Example |
|--------|--------|---------|
| **X-API-Key** (Recommended) | `X-API-Key: your-api-key` | `curl -H "X-API-Key: abc123"` |
| **Bearer Token** | `Authorization: Bearer your-api-key` | `curl -H "Authorization: Bearer abc123"` |
| **API Key** | `Authorization: ApiKey your-api-key` | `curl -H "Authorization: ApiKey abc123"` |
| **Query Parameter** | `?api_key=your-api-key` | `curl "http://localhost:8080/logs?api_key=abc123"` |

### Authentication Errors

| Status Code | Response | Description |
|-------------|----------|-------------|
| **401** | `{"error": "API key required", "hint": "Provide API key in Authorization header, X-API-Key header, or api_key query parameter"}` | No API key provided |
| **401** | `{"error": "Invalid API key"}` | Invalid or expired API key |

## Rate Limiting

- **Limit**: 60 requests per minute per IP address
- **Headers**: All responses include rate limit information

### Rate Limit Headers

```http
X-RateLimit-Limit: 60          # Maximum requests per window
X-RateLimit-Remaining: 45      # Requests remaining in current window
```

### Rate Limit Exceeded

**Status Code**: `429 Too Many Requests`

```json
{
  "error": "Rate limit exceeded",
  "limit": 60,
  "window": "1 minute"
}
```

**Headers**:
```http
Retry-After: 60
```

---

## Health API

### GET /health

System health check endpoint (no authentication required).

**Request**:
```bash
curl http://localhost:8080/health
```

**Response**:
```
Status: 200 OK
Body: OK
```

---

## Logs API

### GET /logs

Retrieve log entries with optional filtering.

**Authentication**: Required

**Query Parameters**:

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `service` | string | Filter by service name | `service=user-service` |
| `level` | string | Filter by log level | `level=ERROR` |
| `message` | string | Search in log messages | `message=login` |
| `start_time` | string (RFC3339) | Start time filter | `start_time=2025-05-26T00:00:00Z` |
| `end_time` | string (RFC3339) | End time filter | `end_time=2025-05-26T23:59:59Z` |
| `limit` | integer | Maximum results (default: 100) | `limit=50` |
| `offset` | integer | Results offset (default: 0) | `offset=100` |

**Request**:
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/logs?service=api-service&level=ERROR&limit=10"
```

**Response**:
```json
[
  {
    "id": 123,
    "service_name": "api-service",
    "log_level": "ERROR",
    "message": "Database connection failed",
    "timestamp": "2025-05-26T10:30:15.123456Z",
    "trace_id": {
      "String": "550e8400-e29b-41d4-a716-446655440000",
      "Valid": true
    },
    "span_id": {
      "String": "6ba7b810-9dad-11d1-80b4-00c04fd430c8", 
      "Valid": true
    },
    "metadata": {
      "error_code": "DB_CONNECTION_TIMEOUT",
      "retry_count": 3
    }
  }
]
```

### POST /logs

Create a new log entry.

**Authentication**: Required

**Request Body**:
```json
{
  "service_name": "user-service",
  "log_level": "INFO",
  "message": "User login successful",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "span_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "metadata": {
    "user_id": "12345",
    "ip_address": "192.168.1.100"
  }
}
```

**Required Fields**:
- `service_name` (string): Name of the service
- `message` (string): Log message content

**Optional Fields**:
- `log_level` (string): One of `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`
- `trace_id` (string): UUID for trace correlation
- `span_id` (string): UUID for span correlation  
- `metadata` (object): Additional structured data

**Request**:
```bash
curl -X POST http://localhost:8080/logs \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "user-service",
    "log_level": "INFO", 
    "message": "User login successful",
    "metadata": {"user_id": "12345"}
  }'
```

**Response**:
```json
{
  "id": 124,
  "service_name": "user-service",
  "log_level": "INFO",
  "message": "User login successful", 
  "timestamp": "2025-05-26T10:35:22.789012Z",
  "trace_id": {
    "String": "",
    "Valid": false
  },
  "span_id": {
    "String": "",
    "Valid": false
  },
  "metadata": {
    "user_id": "12345"
  }
}
```

**Status Codes**:
- `201 Created`: Log entry created successfully
- `400 Bad Request`: Validation error
- `401 Unauthorized`: Authentication required
- `429 Too Many Requests`: Rate limit exceeded

---

## Metrics API

### GET /metrics

Retrieve performance metrics with optional filtering.

**Authentication**: Required

**Query Parameters**:

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `service` | string | Filter by service name | `service=api-gateway` |
| `path` | string | Filter by request path | `path=/api/users` |
| `method` | string | Filter by HTTP method | `method=POST` |
| `min_status` | integer | Minimum status code | `min_status=400` |
| `max_status` | integer | Maximum status code | `max_status=499` |
| `limit` | integer | Maximum results (default: 100) | `limit=50` |
| `offset` | integer | Results offset (default: 0) | `offset=100` |

**Request**:
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/metrics?service=api-service&method=POST&min_status=400"
```

**Response**:
```json
[
  {
    "id": 456,
    "service_name": "api-service",
    "path": "/api/login",
    "method": "POST",
    "status_code": 422,
    "duration_ms": 156.7,
    "source": {
      "language": "go",
      "framework": "gin",
      "version": "1.9.1"
    },
    "environment": "production",
    "timestamp": "2025-05-26T10:28:33.445566Z",
    "request_id": "req-abc-123"
  }
]
```

### POST /metrics

Record a new performance metric.

**Authentication**: Required

**Request Body**:
```json
{
  "service_name": "api-gateway",
  "path": "/api/users",
  "method": "GET", 
  "status_code": 200,
  "duration_ms": 45.7,
  "source": {
    "language": "go",
    "framework": "gin",
    "version": "1.9.1"
  },
  "environment": "production",
  "request_id": "req-xyz-789"
}
```

**Required Fields**:
- `service_name` (string): Name of the service
- `path` (string): Request path or endpoint
- `method` (string): HTTP method (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`)
- `status_code` (integer): HTTP status code (100-599)
- `duration_ms` (number): Request duration in milliseconds
- `source.language` (string): Programming language
- `source.framework` (string): Web framework used
- `source.version` (string): Framework or application version

**Optional Fields**:
- `environment` (string): Deployment environment
- `request_id` (string): Unique request identifier

**Request**:
```bash
curl -X POST http://localhost:8080/metrics \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "api-gateway",
    "path": "/api/users",
    "method": "GET",
    "status_code": 200,
    "duration_ms": 45.7,
    "source": {
      "language": "go", 
      "framework": "gin",
      "version": "1.9.1"
    }
  }'
```

**Response**:
```json
{
  "id": 457,
  "service_name": "api-gateway",
  "path": "/api/users", 
  "method": "GET",
  "status_code": 200,
  "duration_ms": 45.7,
  "source": {
    "language": "go",
    "framework": "gin", 
    "version": "1.9.1"
  },
  "environment": "",
  "timestamp": "2025-05-26T10:40:18.334455Z",
  "request_id": ""
}
```

**Status Codes**:
- `201 Created`: Metric recorded successfully
- `400 Bad Request`: Validation error
- `401 Unauthorized`: Authentication required
- `429 Too Many Requests`: Rate limit exceeded

---

## Traces API

### GET /traces

Retrieve distributed traces with optional filtering.

**Authentication**: Required

**Query Parameters**:

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `service` | string | Filter by service name | `service=api-gateway` |
| `start_time` | string (RFC3339) | Start time filter | `start_time=2025-05-26T00:00:00Z` |
| `end_time` | string (RFC3339) | End time filter | `end_time=2025-05-26T23:59:59Z` |
| `limit` | integer | Maximum results (default: 100) | `limit=50` |
| `offset` | integer | Results offset (default: 0) | `offset=100` |

**Request**:
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/traces?service=api-gateway&limit=10"
```

**Response**:
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "service_name": "api-gateway",
    "start_time": "2025-05-26T10:30:00.123456Z",
    "end_time": {
      "Time": "2025-05-26T10:30:02.456789Z",
      "Valid": true
    },
    "duration_ms": {
      "Float64": 2333.333,
      "Valid": true
    }
  },
  {
    "id": "660f8400-e29b-41d4-a716-446655440001", 
    "service_name": "user-service",
    "start_time": "2025-05-26T10:28:15.789012Z",
    "end_time": {
      "Time": "0001-01-01T00:00:00Z",
      "Valid": false
    },
    "duration_ms": {
      "Float64": 0,
      "Valid": false
    }
  }
]
```

### POST /traces

Create a new distributed trace.

**Authentication**: Required

**Request Body**:
```json
{
  "service_name": "api-gateway"
}
```

**Required Fields**:
- `service_name` (string): Name of the originating service

**Optional Fields**:
- `id` (string): Custom trace ID (UUID generated if not provided)

**Request**:
```bash
curl -X POST http://localhost:8080/traces \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"service_name": "api-gateway"}'
```

**Response**:
```json
{
  "id": "770f8400-e29b-41d4-a716-446655440002",
  "service_name": "api-gateway", 
  "start_time": "2025-05-26T10:45:30.123456Z",
  "end_time": {
    "Time": "0001-01-01T00:00:00Z",
    "Valid": false
  },
  "duration_ms": {
    "Float64": 0,
    "Valid": false
  }
}
```

### POST /traces/{traceId}/end

Mark a trace as completed and calculate duration.

**Authentication**: Required

**Path Parameters**:
- `traceId` (string): UUID of the trace to complete

**Request**:
```bash
curl -X POST http://localhost:8080/traces/770f8400-e29b-41d4-a716-446655440002/end \
  -H "X-API-Key: your-api-key"
```

**Response**:
```json
{
  "id": "770f8400-e29b-41d4-a716-446655440002",
  "service_name": "api-gateway",
  "start_time": "2025-05-26T10:45:30.123456Z", 
  "end_time": {
    "Time": "2025-05-26T10:45:32.789012Z",
    "Valid": true
  },
  "duration_ms": {
    "Float64": 2665.556,
    "Valid": true
  }
}
```

**Status Codes**:
- `200 OK`: Trace completed successfully
- `404 Not Found`: Trace ID not found
- `401 Unauthorized`: Authentication required

---

## Spans API

### GET /traces/{traceId}/spans

Retrieve all spans for a specific trace.

**Authentication**: Required

**Path Parameters**:
- `traceId` (string): UUID of the trace

**Request**:
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/traces/550e8400-e29b-41d4-a716-446655440000/spans"
```

**Response**:
```json
[
  {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "parent_id": "",
    "service": "api-gateway", 
    "operation": "handle_request",
    "start_time": "2025-05-26T10:30:00.123456Z",
    "end_time": "2025-05-26T10:30:00.567890Z",
    "duration_ms": 444.434
  },
  {
    "id": "7ca8c920-adad-22e2-91c5-11d15fe541d9",
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "parent_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "service": "user-service",
    "operation": "authenticate_user", 
    "start_time": "2025-05-26T10:30:00.234567Z",
    "end_time": "2025-05-26T10:30:00.456789Z",
    "duration_ms": 222.222
  }
]
```

### POST /spans

Create a new span within a trace.

**Authentication**: Required

**Request Body**:
```json
{
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "parent_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "service": "database-service",
  "operation": "query_users"
}
```

**Required Fields**:
- `trace_id` (string): UUID of the parent trace
- `service` (string): Name of the service creating the span
- `operation` (string): Description of the operation

**Optional Fields**:
- `parent_id` (string): UUID of parent span (empty for root spans)
- `id` (string): Custom span ID (UUID generated if not provided)

**Request**:
```bash
curl -X POST http://localhost:8080/spans \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "service": "database-service", 
    "operation": "query_users"
  }'
```

**Response**:
```json
{
  "id": "8db9da30-bebe-33f3-a2d6-22e26gf652ea",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "parent_id": "",
  "service": "database-service",
  "operation": "query_users",
  "start_time": "2025-05-26T10:50:15.123456Z", 
  "end_time": "0001-01-01T00:00:00Z",
  "duration_ms": 0
}
```

### POST /spans/{spanId}/end

Mark a span as completed and calculate duration.

**Authentication**: Required

**Path Parameters**:
- `spanId` (string): UUID of the span to complete

**Request**:
```bash
curl -X POST http://localhost:8080/spans/8db9da30-bebe-33f3-a2d6-22e26gf652ea/end \
  -H "X-API-Key: your-api-key"
```

**Response**:
```json
{
  "id": "8db9da30-bebe-33f3-a2d6-22e26gf652ea",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000", 
  "parent_id": "",
  "service": "database-service",
  "operation": "query_users",
  "start_time": "2025-05-26T10:50:15.123456Z",
  "end_time": "2025-05-26T10:50:15.456789Z", 
  "duration_ms": 333.333
}
```

---

## Error Responses

### Common HTTP Status Codes

| Status Code | Description | Common Causes |
|-------------|-------------|---------------|
| **200** | OK | Successful GET request |
| **201** | Created | Successful POST request |
| **400** | Bad Request | Validation error, malformed JSON |
| **401** | Unauthorized | Missing or invalid API key |
| **404** | Not Found | Resource doesn't exist |
| **405** | Method Not Allowed | HTTP method not supported |
| **429** | Too Many Requests | Rate limit exceeded |
| **500** | Internal Server Error | Server-side error |

### Error Response Format

All error responses follow this format:

```json
{
  "error": "Description of the error",
  "hint": "Optional hint for resolution"
}
```

### Validation Errors

**Status Code**: `400 Bad Request`

```json
{
  "error": "service name is required"
}
```

```json
{
  "error": "invalid log level: INVALID"
}
```

```json
{
  "error": "status code must be between 100 and 599"
}
```

## Response Headers

### Standard Headers

All API responses include:

```http
Content-Type: application/json
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
Date: Tue, 26 May 2025 10:30:00 GMT
```

### Rate Limiting Headers

```http
X-RateLimit-Limit: 60           # Maximum requests per window
X-RateLimit-Remaining: 45       # Requests remaining
Retry-After: 60                 # Seconds to wait (when rate limited)
```

## Client Libraries

### Go Client Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Client struct {
    BaseURL string
    APIKey  string
    client  *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        BaseURL: baseURL,
        APIKey:  apiKey,
        client:  &http.Client{},
    }
}

func (c *Client) SendLog(serviceName, level, message string) error {
    logData := map[string]interface{}{
        "service_name": serviceName,
        "log_level":    level,
        "message":      message,
    }
    
    jsonData, _ := json.Marshal(logData)
    req, _ := http.NewRequest("POST", c.BaseURL+"/logs", bytes.NewBuffer(jsonData))
    req.Header.Set("X-API-Key", c.APIKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 201 {
        return fmt.Errorf("API request failed with status %d", resp.StatusCode)
    }
    
    return nil
}
```

### JavaScript Client Example

```javascript
class GoInsightClient {
  constructor(baseURL, apiKey) {
    this.baseURL = baseURL;
    this.apiKey = apiKey;
  }

  async sendLog(serviceName, level, message, metadata = {}) {
    const response = await fetch(`${this.baseURL}/logs`, {
      method: 'POST',
      headers: {
        'X-API-Key': this.apiKey,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        service_name: serviceName,
        log_level: level,
        message: message,
        metadata: metadata
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  }

  async getMetrics(filters = {}) {
    const params = new URLSearchParams(filters);
    const response = await fetch(`${this.baseURL}/metrics?${params}`, {
      headers: {
        'X-API-Key': this.apiKey
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  }
}
```

---

**Last Updated**: May 2025  
**Version**: 1.0.0

For additional examples and integration guides, see the [Usage Guide](usage.md).