# Usage Guide

This guide covers practical usage of Go-Insight for observability in your applications. Learn how to send logs, metrics, and traces, as well as query and analyze your observability data.

## Quick Start

### Authentication Setup

All API calls require authentication. Set up your API key:

```bash
export API_KEY="your-secure-api-key"
export GO_INSIGHT_URL="http://localhost:8080"
```

### Health Check

Verify your Go-Insight instance is running:

```bash
curl $GO_INSIGHT_URL/health
# Returns: OK
```

## Sending Data

### Logs

Send structured log entries to track application events:

#### Basic Log Entry
```bash
curl -X POST $GO_INSIGHT_URL/logs \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "user-service",
    "log_level": "INFO",
    "message": "User login successful",
    "metadata": {
      "user_id": "12345",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0..."
    }
  }'
```

#### Log with Trace Correlation
```bash
curl -X POST $GO_INSIGHT_URL/logs \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "payment-service",
    "log_level": "ERROR",
    "message": "Payment processing failed",
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "span_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "metadata": {
      "error_code": "PAYMENT_DECLINED",
      "amount": 99.99,
      "currency": "USD"
    }
  }'
```

#### Log Levels
Supported log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`

### Metrics

Track performance and business metrics:

#### HTTP Endpoint Metrics
```bash
curl -X POST $GO_INSIGHT_URL/metrics \
  -H "X-API-Key: $API_KEY" \
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
    },
    "environment": "production",
    "request_id": "req-abc-123"
  }'
```

#### Custom Business Metrics
```bash
curl -X POST $GO_INSIGHT_URL/metrics \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "checkout-service",
    "path": "/checkout/complete",
    "method": "POST",
    "status_code": 201,
    "duration_ms": 1250.3,
    "source": {
      "language": "nodejs",
      "framework": "express",
      "version": "4.18.2"
    },
    "environment": "production",
    "request_id": "checkout-xyz-789"
  }'
```

### Distributed Tracing

Track requests across multiple services:

#### Create a Trace
```bash
# Start a new trace
curl -X POST $GO_INSIGHT_URL/traces \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "api-gateway"
  }'

# Response includes trace ID:
# {"id": "550e8400-e29b-41d4-a716-446655440000", "service_name": "api-gateway", "start_time": "2025-05-26T10:30:00Z", ...}
```

#### Create Spans
```bash
# Create spans within the trace
curl -X POST $GO_INSIGHT_URL/spans \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "service": "user-service",
    "operation": "authenticate_user"
  }'

# Child span
curl -X POST $GO_INSIGHT_URL/spans \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "parent_id": "parent-span-id",
    "service": "database-service", 
    "operation": "query_user_table"
  }'
```

#### Complete Trace/Spans
```bash
# End a span
curl -X POST $GO_INSIGHT_URL/spans/{span-id}/end \
  -H "X-API-Key: $API_KEY"

# End the trace
curl -X POST $GO_INSIGHT_URL/traces/{trace-id}/end \
  -H "X-API-Key: $API_KEY"
```

## Querying Data

### Logs

#### Basic Log Retrieval
```bash
# Get all logs (limited to 100 by default)
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/logs"
```

#### Filter by Service
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/logs?service=user-service&limit=10"
```

#### Filter by Log Level
```bash
# Get only error logs
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/logs?level=ERROR&limit=20"
```

#### Filter by Time Range
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/logs?start_time=2025-05-26T00:00:00Z&end_time=2025-05-26T23:59:59Z"
```

#### Search Log Messages
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/logs?message=payment&service=payment-service"
```

#### Pagination
```bash
# Get next page of results
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/logs?limit=50&offset=50"
```

### Metrics

#### Get All Metrics
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/metrics"
```

#### Filter by Service and Path
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/metrics?service=api-gateway&path=/api/users"
```

#### Filter by HTTP Method
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/metrics?method=POST&service=payment-service"
```

#### Filter by Status Code Range
```bash
# Get only error responses (4xx and 5xx)
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/metrics?min_status=400&max_status=599"
```

### Traces

#### Get All Traces
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/traces"
```

#### Filter Traces by Service
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/traces?service=api-gateway"
```

#### Get Spans for a Trace
```bash
curl -H "X-API-Key: $API_KEY" \
  "$GO_INSIGHT_URL/traces/{trace-id}/spans"
```

## Integration Examples

### Application Integration

#### Go Application
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type LogEntry struct {
    ServiceName string            `json:"service_name"`
    LogLevel    string            `json:"log_level"`
    Message     string            `json:"message"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func sendLog(entry LogEntry) error {
    jsonData, _ := json.Marshal(entry)
    
    req, _ := http.NewRequest("POST", "http://localhost:8080/logs", bytes.NewBuffer(jsonData))
    req.Header.Set("X-API-Key", "your-api-key")
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}

// Usage
func main() {
    log := LogEntry{
        ServiceName: "my-go-service",
        LogLevel:    "INFO",
        Message:     "Application started",
        Metadata: map[string]interface{}{
            "version": "1.0.0",
            "port":    8080,
        },
    }
    
    sendLog(log)
}
```

#### Node.js Application
```javascript
const axios = require('axios');

const goInsightClient = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'X-API-Key': 'your-api-key',
    'Content-Type': 'application/json'
  }
});

async function sendLog(serviceName, level, message, metadata = {}) {
  try {
    await goInsightClient.post('/logs', {
      service_name: serviceName,
      log_level: level,
      message: message,
      metadata: metadata
    });
  } catch (error) {
    console.error('Failed to send log:', error.message);
  }
}

async function sendMetric(serviceName, path, method, statusCode, duration) {
  try {
    await goInsightClient.post('/metrics', {
      service_name: serviceName,
      path: path,
      method: method,
      status_code: statusCode,
      duration_ms: duration,
      source: {
        language: 'nodejs',
        framework: 'express',
        version: process.version
      }
    });
  } catch (error) {
    console.error('Failed to send metric:', error.message);
  }
}

// Usage
sendLog('my-node-service', 'INFO', 'User authenticated', { user_id: 12345 });
sendMetric('my-node-service', '/api/login', 'POST', 200, 145.6);
```

#### Python Application
```python
import requests
import json
from datetime import datetime

class GoInsightClient:
    def __init__(self, base_url, api_key):
        self.base_url = base_url
        self.headers = {
            'X-API-Key': api_key,
            'Content-Type': 'application/json'
        }
    
    def send_log(self, service_name, level, message, metadata=None):
        data = {
            'service_name': service_name,
            'log_level': level,
            'message': message
        }
        if metadata:
            data['metadata'] = metadata
            
        try:
            response = requests.post(
                f'{self.base_url}/logs',
                headers=self.headers,
                json=data
            )
            response.raise_for_status()
        except requests.exceptions.RequestException as e:
            print(f'Failed to send log: {e}')
    
    def send_metric(self, service_name, path, method, status_code, duration_ms):
        data = {
            'service_name': service_name,
            'path': path,
            'method': method,
            'status_code': status_code,
            'duration_ms': duration_ms,
            'source': {
                'language': 'python',
                'framework': 'flask',  # or 'django', etc.
                'version': '3.9'
            }
        }
        
        try:
            response = requests.post(
                f'{self.base_url}/metrics',
                headers=self.headers,
                json=data
            )
            response.raise_for_status()
        except requests.exceptions.RequestException as e:
            print(f'Failed to send metric: {e}')

# Usage
client = GoInsightClient('http://localhost:8080', 'your-api-key')
client.send_log('my-python-service', 'INFO', 'Processing started', {'batch_size': 100})
client.send_metric('my-python-service', '/api/process', 'POST', 200, 1250.0)
```

## Monitoring Best Practices

### Log Levels Usage

- **DEBUG**: Detailed information for diagnosing problems
- **INFO**: General information about application flow
- **WARN**: Warning messages for potentially harmful situations
- **ERROR**: Error events that allow the application to continue
- **FATAL**: Very severe error events that may lead to termination

### Structured Logging

Include consistent metadata fields:

```json
{
  "service_name": "user-service",
  "log_level": "INFO",
  "message": "User action performed",
  "metadata": {
    "user_id": "12345",
    "action": "login",
    "ip_address": "192.168.1.100",
    "timestamp": "2025-05-26T10:30:00Z",
    "request_id": "req-abc-123"
  }
}
```

### Trace Correlation

Always include trace and span IDs in logs for request correlation:

```json
{
  "service_name": "payment-service",
  "log_level": "INFO", 
  "message": "Payment processed successfully",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "span_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

### Metric Collection

Track these key metrics for web services:

- **Response time**: Duration of request processing
- **Status codes**: HTTP response status distribution
- **Request volume**: Number of requests per endpoint
- **Error rate**: Percentage of failed requests

## Rate Limiting

Go-Insight implements rate limiting (60 requests per minute per IP). Monitor these headers:

```http
X-RateLimit-Limit: 60          # Maximum requests per window
X-RateLimit-Remaining: 45      # Requests remaining
```

When rate limited (HTTP 429):

```json
{
  "error": "Rate limit exceeded",
  "limit": 60,
  "window": "1 minute"
}
```

## Error Handling

### Common Error Responses

#### Authentication Required (HTTP 401)
```json
{
  "error": "API key required",
  "hint": "Provide API key in Authorization header, X-API-Key header, or api_key query parameter"
}
```

#### Invalid API Key (HTTP 401)
```json
{
  "error": "Invalid API key"
}
```

#### Validation Error (HTTP 400)
```json
{
  "error": "service name is required"
}
```

### Client-Side Error Handling

Always implement proper error handling in your applications:

```javascript
try {
  await goInsightClient.post('/logs', logData);
} catch (error) {
  if (error.response.status === 429) {
    // Rate limited - implement backoff
    console.log('Rate limited, retrying after delay...');
  } else if (error.response.status === 401) {
    // Authentication issue
    console.error('Invalid API key');
  } else {
    // Other errors
    console.error('Failed to send log:', error.message);
  }
}
```

## Visualizing Logs

After starting the Docker stack, Grafana runs on `http://localhost:3000` with default credentials `admin/admin`. The provided dashboard shows recent log entries from the `logs` table. You can customize queries or create new dashboards to fit your needs.

## Performance Tips

1. **Batch Operations**: Group multiple log entries when possible
2. **Async Sending**: Send observability data asynchronously to avoid blocking main application flow
3. **Error Handling**: Implement retry logic with exponential backoff
4. **Connection Pooling**: Reuse HTTP connections for better performance
5. **Rate Limit Awareness**: Monitor rate limit headers and implement client-side throttling

## Support

For usage questions or integration help:

- Check the [API Reference](api.md) for detailed endpoint documentation
- Review [Security Guide](security.md) for authentication setup
- See [Performance Guide](performance.md) for optimization tips

---

**Last Updated**: May 2025  
**Version**: 1.0.0