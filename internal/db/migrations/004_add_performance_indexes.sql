-- Index for service-based log queries (most common filter)
CREATE INDEX IF NOT EXISTS idx_logs_service_timestamp 
ON logs(service_name, timestamp DESC);

-- Index for log level filtering
CREATE INDEX IF NOT EXISTS idx_logs_level_timestamp 
ON logs(log_level, timestamp DESC);

-- Index for trace correlation (joining logs to traces)
CREATE INDEX IF NOT EXISTS idx_logs_trace_id 
ON logs(trace_id) WHERE trace_id IS NOT NULL;

-- Index for service-based metrics queries
CREATE INDEX IF NOT EXISTS idx_metrics_service_timestamp 
ON metrics(service_name, timestamp DESC);

-- Index for path-based performance analysis
CREATE INDEX IF NOT EXISTS idx_metrics_path_timestamp 
ON metrics(path, timestamp DESC);

-- Index for method + status code analysis
CREATE INDEX IF NOT EXISTS idx_metrics_method_status 
ON metrics(method, status_code);

-- Composite index for error rate analysis
CREATE INDEX IF NOT EXISTS idx_metrics_service_status_time
ON metrics(service_name, status_code, timestamp DESC);

-- Index for service-based trace queries
CREATE INDEX IF NOT EXISTS idx_traces_service_start_time 
ON traces(service_name, start_time DESC);

-- Index for trace duration analysis
CREATE INDEX IF NOT EXISTS idx_traces_duration 
ON traces(duration_ms) WHERE duration_ms IS NOT NULL;

-- Index for active traces (end_time is NULL)
CREATE INDEX IF NOT EXISTS idx_traces_active 
ON traces(start_time DESC) WHERE end_time IS NULL;

-- Index for finding spans by trace (most common query)
CREATE INDEX IF NOT EXISTS idx_spans_trace_id_start_time 
ON spans(trace_id, start_time ASC);

-- Index for parent-child span relationships
CREATE INDEX IF NOT EXISTS idx_spans_parent_id 
ON spans(parent_id) WHERE parent_id IS NOT NULL;

-- Index for service-based span analysis
CREATE INDEX IF NOT EXISTS idx_spans_service_operation 
ON spans(service, operation);