CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    log_level TEXT CHECK (log_level IN ('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL')),
    message TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    trace_id UUID,
    span_id UUID,
    metadata JSONB
);