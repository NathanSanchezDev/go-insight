CREATE TABLE IF NOT EXISTS traces (
    id TEXT PRIMARY KEY,
    service_name TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_ms FLOAT
);

CREATE TABLE IF NOT EXISTS spans (
    id TEXT PRIMARY KEY,
    trace_id TEXT NOT NULL,
    parent_id TEXT,
    service TEXT NOT NULL,
    operation TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_ms FLOAT,
    FOREIGN KEY(trace_id) REFERENCES traces(id)
);
