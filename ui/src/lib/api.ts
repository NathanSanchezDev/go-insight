// ui/src/lib/api.ts

const API_BASE_URL = '/api';

// Types for our API responses
export interface Log {
  id: number;
  service_name: string;
  log_level: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'FATAL';
  message: string;
  timestamp: string;
  trace_id?: { String: string; Valid: boolean };
  span_id?: { String: string; Valid: boolean };
  metadata?: Record<string, any>;
}

export interface Metric {
  id: number;
  service_name: string;
  path: string;
  method: string;
  status_code: number;
  duration_ms: number;
  source: {
    language: string;
    framework: string;
    version: string;
  };
  environment?: string;
  timestamp: string;
  request_id?: string;
}

export interface Trace {
  id: string;
  service_name: string;
  start_time: string;
  end_time?: { Time: string; Valid: boolean };
  duration_ms?: { Float64: number; Valid: boolean };
}

export interface Span {
  id: string;
  trace_id: string;
  parent_id?: string;
  service: string;
  operation: string;
  start_time: string;
  end_time?: string;
  duration_ms?: number;
}

// API query parameters
export interface LogsQuery {
  service?: string;
  level?: string;
  message?: string;
  start_time?: string;
  end_time?: string;
  limit?: number;
  offset?: number;
}

export interface MetricsQuery {
  service?: string;
  path?: string;
  method?: string;
  min_status?: number;
  max_status?: number;
  limit?: number;
  offset?: number;
}

export interface TracesQuery {
  service?: string;
  start_time?: string;
  end_time?: string;
  limit?: number;
  offset?: number;
}

class ApiClient {
  private baseUrl: string;
  private apiKey: string | null = null;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
    // For now, we'll handle auth later
    this.apiKey = 'your-secure-api-key-here'; // TODO: Make this configurable
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers as Record<string, string>,
    };

    if (this.apiKey) {
      headers['X-API-Key'] = this.apiKey;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      throw new Error(`API request failed: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  private buildQueryString(params: Record<string, any>): string {
    const searchParams = new URLSearchParams();
    
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        searchParams.append(key, value.toString());
      }
    });

    const queryString = searchParams.toString();
    return queryString ? `?${queryString}` : '';
  }

  // Health check (no auth required)
  async getHealth(): Promise<string> {
    const response = await fetch(`${this.baseUrl}/health`);
    return response.text();
  }

  // Logs API
  async getLogs(query: LogsQuery = {}): Promise<Log[]> {
    const queryString = this.buildQueryString(query);
    return this.request<Log[]>(`/logs${queryString}`);
  }

  async createLog(log: Partial<Log>): Promise<Log> {
    return this.request<Log>('/logs', {
      method: 'POST',
      body: JSON.stringify(log),
    });
  }

  // Metrics API
  async getMetrics(query: MetricsQuery = {}): Promise<Metric[]> {
    const queryString = this.buildQueryString(query);
    return this.request<Metric[]>(`/metrics${queryString}`);
  }

  async createMetric(metric: Partial<Metric>): Promise<Metric> {
    return this.request<Metric>('/metrics', {
      method: 'POST',
      body: JSON.stringify(metric),
    });
  }

  // Traces API
  async getTraces(query: TracesQuery = {}): Promise<Trace[]> {
    const queryString = this.buildQueryString(query);
    return this.request<Trace[]>(`/traces${queryString}`);
  }

  async createTrace(trace: Partial<Trace>): Promise<Trace> {
    return this.request<Trace>('/traces', {
      method: 'POST',
      body: JSON.stringify(trace),
    });
  }

  async endTrace(traceId: string): Promise<Trace> {
    return this.request<Trace>(`/traces/${traceId}/end`, {
      method: 'POST',
    });
  }

  // Spans API
  async getSpans(traceId: string): Promise<Span[]> {
    return this.request<Span[]>(`/traces/${traceId}/spans`);
  }

  async createSpan(span: Partial<Span>): Promise<Span> {
    return this.request<Span>('/spans', {
      method: 'POST',
      body: JSON.stringify(span),
    });
  }

  async endSpan(spanId: string): Promise<Span> {
    return this.request<Span>(`/spans/${spanId}/end`, {
      method: 'POST',
    });
  }

  // Dashboard summary methods (derived from existing data)
  async getDashboardStats(): Promise<{
    totalLogs: number;
    errorCount: number;
    avgResponseTime: number;
    activeTraces: number;
    topServices: string[];
  }> {
    // Get recent data for dashboard stats
    const [logs, metrics, traces] = await Promise.all([
      this.getLogs({ limit: 1000 }),
      this.getMetrics({ limit: 1000 }),
      this.getTraces({ limit: 100 }),
    ]);

    const errorCount = logs.filter(log => log.log_level === 'ERROR').length;
    const avgResponseTime = metrics.length > 0 
      ? metrics.reduce((sum, m) => sum + m.duration_ms, 0) / metrics.length 
      : 0;
    const activeTraces = traces.filter(t => !t.end_time?.Valid).length;
    
    const serviceSet = new Set([
      ...logs.map(l => l.service_name),
      ...metrics.map(m => m.service_name),
      ...traces.map(t => t.service_name),
    ]);
    const topServices = Array.from(serviceSet).slice(0, 5);

    return {
      totalLogs: logs.length,
      errorCount,
      avgResponseTime: Math.round(avgResponseTime * 100) / 100,
      activeTraces,
      topServices,
    };
  }
}

export const apiClient = new ApiClient();