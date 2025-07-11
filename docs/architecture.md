# Architecture Overview

Go-Insight follows a modular, layered architecture designed for production observability workloads. This document covers the system design, component interactions, data flow, and architectural decisions.

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Go-Insight Platform                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Clients   │  │   Load      │  │  Monitoring │        │
│  │Applications │  │ Balancer    │  │   Systems   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│         │                │                │                │
│         └────────────────┼────────────────┘                │
│                          │                                 │
├─────────────────────────────────────────────────────────────┤
│                    API Gateway Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │Rate Limiting│  │Authentication│  │   Request   │        │
│  │ Middleware  │  │  Middleware  │  │  Logging    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                   Application Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Logs     │  │   Metrics   │  │   Traces    │        │
│  │   Handler   │  │   Handler   │  │   Handler   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                   Business Logic Layer                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Data     │  │ Validation  │  │ Correlation │        │
│  │ Processing  │  │   Engine    │  │   Engine    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                   Data Access Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Connection  │  │   Query     │  │ Transaction │        │
│  │    Pool     │  │ Optimizer   │  │  Manager    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                      Storage Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ PostgreSQL  │  │   Indexes   │  │ Migrations  │        │
│  │  Database   │  │   Strategy  │  │   System    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. API Gateway Layer

#### Rate Limiting Middleware
**Purpose**: Prevent API abuse and ensure fair resource allocation

**Implementation**:
```go
// Per-IP rate limiting with token bucket algorithm
type RateLimiter struct {
    clients map[string]*bucketInfo
    mutex   sync.RWMutex
    limit   int
    window  time.Duration
}
```

**Features**:
- Per-IP tracking with separate buckets
- Memory-efficient with automatic cleanup
- Thread-safe concurrent access
- Real-time header feedback (`X-RateLimit-Remaining`)

#### Authentication Middleware
**Purpose**: Secure API access with flexible authentication methods

**Implementation**:
```go
// Multi-method authentication support
func extractAPIKey(r *http.Request) string {
    // 1. Authorization Bearer
    // 2. Authorization ApiKey  
    // 3. X-API-Key header
    // 4. Query parameter
}
```

**Features**:
- Multiple authentication methods for client flexibility
- Public/protected endpoint separation
- Clear error messages with helpful hints
- Environment-based configuration

#### Request Logging Middleware
**Purpose**: Audit trail and performance monitoring

**Implementation**:
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
    })
}
```

### 2. Application Layer

#### Logs Handler
**Responsibilities**:
- Log ingestion with validation
- Filtering and querying with pagination
- Trace correlation via trace_id/span_id
- Structured metadata support

**Key Features**:
- Support for standard log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- Full-text search capabilities
- Time-range filtering with optimized queries
- JSON metadata storage with flexible schema

#### Metrics Handler  
**Responsibilities**:
- Performance metrics collection
- HTTP endpoint monitoring
- Custom business metrics
- Source attribution (language, framework, version)

**Key Features**:
- HTTP method and status code tracking
- Response time measurement
- Environment-based metrics separation
- Request correlation with unique IDs

#### Traces Handler
**Responsibilities**:
- Distributed trace lifecycle management
- Span creation and completion
- Parent-child relationship modeling
- Duration calculation and tracking

**Key Features**:
- UUID-based trace and span identification
- Hierarchical span relationships
- Automatic duration calculation
- Cross-service request tracking

### 3. Data Access Layer

#### Connection Pool Management
**Implementation**:
```go
// PostgreSQL connection pooling
func InitDB() {
    DB, err = sql.Open("pgx", dsn)
    DB.SetMaxOpenConns(25)      // Maximum connections
    DB.SetMaxIdleConns(10)      // Idle connection pool
    DB.SetConnMaxLifetime(time.Hour)  // Connection lifetime
}
```

**Benefits**:
- Efficient resource utilization
- Prevents connection exhaustion
- Automatic connection lifecycle management
- Performance optimization through reuse

#### Query Optimization
**Strategic Indexing**:
```sql
-- Service-based queries (most common pattern)
CREATE INDEX idx_logs_service_timestamp ON logs(service_name, timestamp DESC);

-- Composite indexes for complex filtering
CREATE INDEX idx_metrics_service_status_time 
ON metrics(service_name, status_code, timestamp DESC);

-- Trace correlation
CREATE INDEX idx_logs_trace_id ON logs(trace_id) WHERE trace_id IS NOT NULL;
```

**Performance Results**:
- 10-100x query performance improvement
- Sub-10ms response times for filtered queries
- Efficient pagination and sorting

## Data Flow Architecture

### 1. Ingestion Flow

```
Client Application
       │
       ▼
┌─────────────┐
│Rate Limiting│ ── Check: 60 req/min per IP
└─────────────┘
       │
       ▼
┌─────────────┐
│Authentication│ ── Validate: API key required
└─────────────┘
       │
       ▼
┌─────────────┐
│  Validation │ ── Check: Required fields, data types
└─────────────┘
       │
       ▼
┌─────────────┐
│  Processing │ ── Add: Timestamps, UUIDs
└─────────────┘
       │
       ▼
┌─────────────┐
│   Storage   │ ── Store: PostgreSQL with indexes
└─────────────┘
       │
       ▼
┌─────────────┐
│   Response  │ ── Return: Created entity with ID
└─────────────┘
```

### 2. Query Flow

```
Client Query Request
       │
       ▼
┌─────────────┐
│Rate Limiting│ ── Check: Request allowance
└─────────────┘
       │
       ▼
┌─────────────┐
│Authentication│ ── Validate: API credentials
└─────────────┘
       │
       ▼
┌─────────────┐
│Query Builder│ ── Parse: Filters, pagination, sorting
└─────────────┘
       │
       ▼
┌─────────────┐
│Index Lookup │ ── Optimize: Use strategic indexes
└─────────────┘
       │
       ▼
┌─────────────┐
│Result Set   │ ── Apply: Limits, offsets, ordering
└─────────────┘
       │
       ▼
┌─────────────┐
│Serialization│ ── Format: JSON response
└─────────────┘
```

### 3. Distributed Tracing Flow

```
Request Start
       │
       ▼
┌─────────────┐
│Create Trace │ ── Generate: UUID, start timestamp
└─────────────┘
       │
       ▼
┌─────────────┐
│Create Spans │ ── Track: Service operations
└─────────────┘
       │
       ▼
┌─────────────┐
│ Correlation │ ── Link: Logs with trace_id/span_id
└─────────────┘
       │
       ▼
┌─────────────┐
│End Spans    │ ── Calculate: Duration, mark complete
└─────────────┘
       │
       ▼
┌─────────────┐
│ End Trace   │ ── Finalize: Total duration
└─────────────┘
```

## Database Schema Design

### Entity Relationship Diagram

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│    LOGS     │       │   METRICS   │       │   TRACES    │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ id (PK)     │       │ id (PK)     │       │ id (PK)     │
│ service_name│       │ service_name│       │ service_name│
│ log_level   │       │ path        │       │ start_time  │
│ message     │       │ method      │       │ end_time    │
│ timestamp   │       │ status_code │       │ duration_ms │
│ trace_id(FK)│────┐  │ duration    │       └─────────────┘
│ span_id(FK) │    │  │ language    │              │
│ metadata    │    │  │ framework   │              │
└─────────────┘    │  │ version     │              │
                   │  │ environment │              │
                   │  │ timestamp   │              │
                   │  │ request_id  │              │
                   │  └─────────────┘              │
                   │                               │
                   │  ┌─────────────┐              │
                   └──│    SPANS    │──────────────┘
                      ├─────────────┤
                      │ id (PK)     │
                      │ trace_id(FK)│
                      │ parent_id   │
                      │ service     │
                      │ operation   │
                      │ start_time  │
                      │ end_time    │
                      │ duration_ms │
                      └─────────────┘
```

### Schema Design Principles

1. **Normalization**: Separate concerns while maintaining query efficiency
2. **Indexing Strategy**: Composite indexes for common query patterns
3. **Data Types**: Appropriate types for performance (UUID, TIMESTAMP, JSONB)
4. **Constraints**: Foreign keys for referential integrity
5. **Null Handling**: Proper nullable types for optional relationships

## Security Architecture

### Defense in Depth

```
┌─────────────────────────────────────────────────────────┐
│                    Security Layers                     │
├─────────────────────────────────────────────────────────┤
│ Layer 1: Network (Load Balancer, Firewall)            │
├─────────────────────────────────────────────────────────┤
│ Layer 2: Rate Limiting (Per-IP throttling)            │
├─────────────────────────────────────────────────────────┤
│ Layer 3: Authentication (API key validation)          │
├─────────────────────────────────────────────────────────┤
│ Layer 4: Authorization (Endpoint access control)      │
├─────────────────────────────────────────────────────────┤
│ Layer 5: Input Validation (Data sanitization)         │
├─────────────────────────────────────────────────────────┤
│ Layer 6: Database Security (Parameterized queries)    │
└─────────────────────────────────────────────────────────┘
```

### Security Components

#### 1. Authentication System
- **Multi-method support**: Headers, Bearer tokens, query parameters
- **Environment-based keys**: Different keys per deployment environment
- **Clear error responses**: Helpful hints without exposing internals

#### 2. Rate Limiting System
- **Per-IP tracking**: Separate limits for different client IPs
- **Token bucket algorithm**: Smooth rate limiting with burst capacity
- **Automatic cleanup**: Memory leak prevention
- **Real-time feedback**: Rate limit headers for client awareness

#### 3. Input Validation
- **Type validation**: Ensure correct data types for all fields
- **Range validation**: Status codes, timestamps, required fields
- **Sanitization**: Prevent injection attacks in log messages
- **Size limits**: Prevent oversized payloads

## Performance Architecture

### Performance Optimization Strategies

#### 1. Database Layer
```sql
-- Strategic indexing for common query patterns
CREATE INDEX idx_logs_service_timestamp ON logs(service_name, timestamp DESC);
CREATE INDEX idx_metrics_service_status_time ON metrics(service_name, status_code, timestamp DESC);
CREATE INDEX idx_traces_service_start_time ON traces(service_name, start_time DESC);
```

**Results**:
- 10-100x query performance improvement
- Sub-10ms response times
- Efficient memory usage

#### 2. Application Layer
- **Connection pooling**: Reuse database connections
- **Prepared statements**: Query plan caching
- **Efficient serialization**: Minimal JSON processing overhead
- **Memory management**: Automatic cleanup and garbage collection

#### 3. Concurrency Design
- **Thread-safe middleware**: Proper synchronization primitives
- **Non-blocking operations**: Asynchronous processing where possible
- **Resource isolation**: Separate rate limiting buckets per IP
- **Graceful degradation**: Continue operation under high load

### Performance Monitoring

#### Built-in Metrics
- Request processing time
- Database query duration  
- Rate limit utilization
- Connection pool usage
- Memory allocation patterns

#### Performance Headers
```http
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
```

## Deployment Architecture

### Single Instance Deployment

```
┌─────────────────────────────────────────┐
│              Load Balancer              │
│            (nginx/Apache)               │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│             Go-Insight                  │
│         (Single Instance)               │
│  ┌─────────────┐ ┌─────────────────────┐│
│  │   API       │ │    Middleware       ││
│  │  Handlers   │ │  - Authentication   ││
│  │             │ │  - Rate Limiting    ││
│  │             │ │  - Logging          ││
│  └─────────────┘ └─────────────────────┘│
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            PostgreSQL                   │
│          (Primary Database)             │
└─────────────────────────────────────────┘
```

### High Availability Deployment (Future)

```
┌─────────────────────────────────────────┐
│              Load Balancer              │
│         (with Health Checks)            │
└────┬─────────────────────────────┬──────┘
     │                             │
┌────▼─────┐                 ┌─────▼────┐
│Go-Insight│                 │Go-Insight│
│Instance 1│                 │Instance 2│
└────┬─────┘                 └─────┬────┘
     │                             │
     └─────────────┬─────────────────┘
                   │
┌─────────────────▼───────────────────────┐
│            PostgreSQL Cluster           │
│  ┌─────────────┐ ┌─────────────────────┐│
│  │   Primary   │ │     Read Replicas   ││
│  │  (Write)    │ │     (Read Only)     ││
│  └─────────────┘ └─────────────────────┘│
└─────────────────────────────────────────┘
```

## Scalability Considerations

### Current Scalability Features
- **Thread-safe design**: Handles concurrent requests safely
- **Connection pooling**: Efficient database resource usage
- **Strategic indexing**: Maintains performance as data grows
- **Rate limiting**: Prevents resource exhaustion

### Future Scalability Enhancements
- **Horizontal scaling**: Multiple Go-Insight instances
- **Database sharding**: Distribute data across multiple databases
- **Caching layer**: Redis for frequently accessed data
- **Async processing**: Background jobs for heavy operations

## Monitoring and Observability

### Internal Monitoring
Go-Insight monitors its own performance:

```go
// Request timing
start := time.Now()
next.ServeHTTP(w, r)
duration := time.Since(start)
log.Printf("%s %s %s", r.Method, r.RequestURI, duration)
```

### Health Checks
- **Database connectivity**: Verify PostgreSQL connection
- **Application health**: Memory usage, goroutine count
- **System resources**: CPU, memory, disk space

### Logging Strategy
- **Structured logging**: Consistent log format
- **Security events**: Authentication failures, rate limiting
- **Performance events**: Slow queries, high memory usage
- **Error tracking**: Application errors and stack traces

## Technology Stack

### Core Technologies
- **Language**: Go 1.23+
- **Database**: PostgreSQL 15+
- **HTTP Router**: Gorilla Mux
- **Database Driver**: pgx/v5

### Dependencies
```go
// go.mod
require (
    github.com/gorilla/mux v1.8.1
    github.com/jackc/pgx/v5 v5.7.2
    github.com/joho/godotenv v1.5.1
    github.com/google/uuid v1.6.0
)
```

### Development Tools
- **Docker**: Containerized PostgreSQL for development
- **Migration System**: SQL-based schema versioning
- **Environment Management**: .env file configuration

## Design Decisions

### 1. Database Choice: PostgreSQL
**Why PostgreSQL**:
- ACID compliance for data integrity
- Advanced indexing capabilities (GIN, composite indexes)
- JSONB support for flexible metadata
- Excellent performance for read-heavy workloads
- Strong ecosystem and tooling

**Alternatives Considered**:
- **ClickHouse**: Better for analytics, but less flexible for mixed workloads
- **InfluxDB**: Time-series focused, but limited relational capabilities
- **MongoDB**: Document store, but less predictable performance

### 2. Programming Language: Go
**Why Go**:
- Excellent concurrency primitives (goroutines, channels)
- Strong standard library for HTTP and database operations
- Fast compilation and deployment
- Memory efficient with automatic garbage collection
- Strong typing system prevents many runtime errors

**Alternatives Considered**:
- **Rust**: Better performance, but steeper learning curve
- **Java**: Enterprise ecosystem, but higher memory usage
- **Node.js**: JavaScript ecosystem, but single-threaded limitations

### 3. Architecture Pattern: Layered Architecture
**Why Layered Architecture**:
- Clear separation of concerns
- Easy to test individual components
- Flexible for future enhancements
- Well-understood pattern for team development

**Alternatives Considered**:
- **Microservices**: More complex for current scope
- **Hexagonal Architecture**: Over-engineered for current needs
- **Event-Driven**: Adds complexity without clear benefit

### 4. Security Model: API Key Authentication
**Why API Keys**:
- Simple to implement and understand
- Sufficient for current security requirements
- Easy to rotate and manage
- Low overhead for high-throughput APIs

**Current Enhancements**:
- JWT tokens for stateless authentication
- Role-based access control (RBAC)
**Future**: OAuth2 integration for enterprise environments

### 5. Log Visualization UI
**Decision**: Integrate Grafana for log dashboards instead of building a custom UI in the short term.

**Why Grafana**:
- Mature visualization platform with built-in PostgreSQL support
- Ready-made dashboarding capabilities reduce development time
- Easily extensible for metrics and traces

**Alternatives Considered**:
- **Custom React UI**: Full control over design but requires significant effort
- **Kibana**: Another log UI option, but heavier stack and less integrated with existing plans

Grafana can be deployed alongside Go-Insight via Docker Compose and connects directly to the PostgreSQL database. A custom UI remains possible in the future if specialized features are needed.

## Testing Strategy

### Unit Testing
```go
// Example test structure
func TestRateLimitMiddleware(t *testing.T) {
    // Test rate limiting logic
    // Verify thread safety
    // Check cleanup functionality
}
```

### Integration Testing
- Database operations with test database
- API endpoint testing with authentication
- End-to-end trace correlation testing

### Performance Testing
- Load testing with concurrent requests
- Database performance under various data sizes
- Memory usage profiling

### Security Testing
- Authentication bypass attempts
- Rate limit circumvention testing
- Input validation with malicious payloads

## Future Architecture Evolution

### Phase 2: Enhanced Production Features
- **Input validation layer**: Request sanitization and size limits
- **Background job processing**: Async data processing
- **Internal metrics collection**: Self-monitoring capabilities

### Phase 3: User Interface Layer
- **Web frontend**: React-based dashboard
- **Real-time updates**: WebSocket connections
- **Visualization engine**: Time-series charts and graphs

### Phase 4: Advanced Analytics
- **Stream processing**: Real-time data analysis
- **Machine learning**: Anomaly detection
- **Advanced querying**: Custom query language

### Phase 5: Enterprise Features
- **Multi-tenancy**: Data isolation per organization
- **High availability**: Distributed deployment
- **Advanced security**: RBAC, audit logging

---

**Last Updated**: May 2025  
**Version**: 1.0.0