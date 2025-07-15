// ui/src/hooks/use-dashboard-data.ts
import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"

export function useDashboardStats() {
  return useQuery({
    queryKey: ["dashboard-stats"],
    queryFn: () => apiClient.getDashboardStats(),
    refetchInterval: 30000, // Refetch every 30 seconds
  })
}

export function useRecentLogs(limit = 5) {
  return useQuery({
    queryKey: ["recent-logs", limit],
    queryFn: () => apiClient.getLogs({ limit }),
    refetchInterval: 15000, // Refetch every 15 seconds
  })
}

export function useSystemHealth() {
  return useQuery({
    queryKey: ["system-health"],
    queryFn: async () => {
      // Get health status and recent metrics for service status
      const [health, metrics] = await Promise.all([
        apiClient.getHealth(),
        apiClient.getMetrics({ limit: 100 })
      ])

      // Group metrics by service to determine health
      const serviceMetrics = metrics.reduce((acc, metric) => {
        if (!acc[metric.service_name]) {
          acc[metric.service_name] = []
        }
        acc[metric.service_name].push(metric)
        return acc
      }, {} as Record<string, typeof metrics>)

      // Determine service health based on recent error rates
      const services = Object.entries(serviceMetrics).map(([serviceName, serviceMetrics]) => {
        const recentMetrics = serviceMetrics.slice(0, 10) // Last 10 requests
        const errorCount = recentMetrics.filter(m => m.status_code >= 500).length
        const errorRate = recentMetrics.length > 0 ? errorCount / recentMetrics.length : 0
        
        let status: "healthy" | "warning" | "error"
        if (errorRate === 0) {
          status = "healthy"
        } else if (errorRate < 0.1) {
          status = "warning"
        } else {
          status = "error"
        }

        return {
          name: serviceName,
          status,
          errorRate: Math.round(errorRate * 100),
          requestCount: recentMetrics.length
        }
      })

      return { health, services }
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  })
}