CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    duration DOUBLE PRECISION NOT NULL,
    language TEXT NOT NULL,
    framework TEXT NOT NULL,
    version TEXT NOT NULL,
    environment TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    request_id TEXT
);
