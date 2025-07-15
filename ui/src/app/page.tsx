// ui/src/app/page.tsx
"use client"

import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Activity, AlertTriangle, Clock, Zap, RefreshCw } from "lucide-react"
import { useDashboardStats, useRecentLogs, useSystemHealth } from "@/hooks/use-dashboard-data"
import { Button } from "@/components/ui/button"

function StatsCards() {
  const { data: stats, isLoading, refetch, isFetching } = useDashboardStats()

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Card key={i}>
            <CardHeader className="pb-2">
              <Skeleton className="h-4 w-20" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-16 mb-1" />
              <Skeleton className="h-3 w-24" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  if (!stats) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <p className="text-sm text-muted-foreground">Failed to load stats</p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Logs</CardTitle>
          <Activity className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.totalLogs.toLocaleString()}</div>
          <p className="text-xs text-muted-foreground">
            {isFetching ? "Updating..." : "Real-time data"}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Error Count</CardTitle>
          <AlertTriangle className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.errorCount}</div>
          <p className="text-xs text-muted-foreground">
            Error-level log entries
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Avg Response Time</CardTitle>
          <Clock className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.avgResponseTime}ms</div>
          <p className="text-xs text-muted-foreground">
            Across all services
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Active Traces</CardTitle>
          <Zap className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.activeTraces}</div>
          <p className="text-xs text-muted-foreground">
            Currently in progress
          </p>
        </CardContent>
      </Card>
    </div>
  )
}

function RecentLogsCard() {
  const { data: logs, isLoading, refetch, isFetching } = useRecentLogs(5)

  const getLogLevelBadgeVariant = (level: string) => {
    switch (level.toUpperCase()) {
      case 'ERROR':
      case 'FATAL':
        return 'destructive'
      case 'WARN':
        return 'secondary'
      default:
        return 'outline'
    }
  }

  const formatTimeAgo = (timestamp: string) => {
    const now = new Date()
    const logTime = new Date(timestamp)
    const diffInSeconds = Math.floor((now.getTime() - logTime.getTime()) / 1000)
    
    if (diffInSeconds < 60) return `${diffInSeconds}s ago`
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`
    return `${Math.floor(diffInSeconds / 3600)}h ago`
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Recent Logs</CardTitle>
            <CardDescription>
              Latest log entries from your services
            </CardDescription>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => refetch()}
            disabled={isFetching}
          >
            <RefreshCw className={`h-4 w-4 ${isFetching ? 'animate-spin' : ''}`} />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-2">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="flex items-center justify-between p-2">
                <div className="flex items-center gap-2">
                  <Skeleton className="h-5 w-12" />
                  <Skeleton className="h-4 w-40" />
                </div>
                <Skeleton className="h-3 w-12" />
              </div>
            ))}
          </div>
        ) : logs && logs.length > 0 ? (
          <div className="space-y-2">
            {logs.map((log) => (
              <div key={log.id} className="flex items-center justify-between p-2 rounded bg-muted/50">
                <div className="flex items-center gap-2">
                  <Badge variant={getLogLevelBadgeVariant(log.log_level)}>
                    {log.log_level || 'INFO'}
                  </Badge>
                  <span className="text-sm truncate max-w-[200px]" title={log.message}>
                    {log.message}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    [{log.service_name}]
                  </span>
                </div>
                <span className="text-xs text-muted-foreground">
                  {formatTimeAgo(log.timestamp)}
                </span>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">No logs available</p>
        )}
      </CardContent>
    </Card>
  )
}

function ServiceHealthCard() {
  const { data: healthData, isLoading, refetch, isFetching } = useSystemHealth()

  const getStatusBadgeProps = (status: string) => {
    switch (status) {
      case 'healthy':
        return { className: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100" }
      case 'warning':
        return { className: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-100" }
      case 'error':
        return { className: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-100" }
      default:
        return { variant: "outline" as const }
    }
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Service Health</CardTitle>
            <CardDescription>
              Current status of your services
            </CardDescription>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => refetch()}
            disabled={isFetching}
          >
            <RefreshCw className={`h-4 w-4 ${isFetching ? 'animate-spin' : ''}`} />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="flex items-center justify-between p-2">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-5 w-16" />
              </div>
            ))}
          </div>
        ) : healthData?.services && healthData.services.length > 0 ? (
          <div className="space-y-2">
            {healthData.services.map((service) => (
              <div key={service.name} className="flex items-center justify-between p-2 rounded bg-muted/50">
                <span className="text-sm font-medium">{service.name}</span>
                <div className="flex items-center gap-2">
                  <span className="text-xs text-muted-foreground">
                    {service.errorRate}% errors
                  </span>
                  <Badge {...getStatusBadgeProps(service.status)}>
                    {service.status}
                  </Badge>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">No services detected</p>
        )}
      </CardContent>
    </Card>
  )
}

export default function DashboardHome() {
  const { data: healthData } = useSystemHealth()
  
  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground">
              Monitor your distributed applications in real-time
            </p>
          </div>
          <Badge 
            variant="outline" 
            className={
              healthData?.health === "OK" 
                ? "text-green-600 border-green-600" 
                : "text-red-600 border-red-600"
            }
          >
            <Activity className="mr-1 h-3 w-3" />
            {healthData?.health === "OK" ? "System Healthy" : "System Issues"}
          </Badge>
        </div>

        {/* Stats Cards */}
        <StatsCards />

        {/* Recent Activity */}
        <div className="grid gap-4 md:grid-cols-2">
          <RecentLogsCard />
          <ServiceHealthCard />
        </div>
      </div>
    </DashboardLayout>
  )
}