// ui/src/components/dashboard-layout.tsx
"use client"

import * as React from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { 
  BarChart3, 
  FileText, 
  Activity, 
  Settings, 
  Moon, 
  Sun, 
  Monitor,
  Home,
  Network
} from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { 
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarProvider,
  SidebarRail,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarFooter
} from "@/components/ui/sidebar"
import { useTheme } from "./theme-provider"

const navigationItems = [
  {
    title: "Overview",
    url: "/",
    icon: Home,
  },
  {
    title: "Logs",
    url: "/logs", 
    icon: FileText,
  },
  {
    title: "Metrics",
    url: "/metrics",
    icon: BarChart3,
  },
  {
    title: "Traces",
    url: "/traces",
    icon: Network,
  },
  {
    title: "System Health",
    url: "/health",
    icon: Activity,
  },
]

function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => {
        setTheme(
          theme === "light" 
            ? "dark" 
            : theme === "dark" 
            ? "system" 
            : "light"
        )
      }}
      className="w-9 px-0"
    >
      <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
      <Monitor className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all" />
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}

function AppSidebar() {
  const pathname = usePathname()

  return (
    <Sidebar>
      <SidebarHeader className="border-b border-sidebar-border">
        <div className="flex items-center gap-2 px-4 py-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Activity className="h-4 w-4" />
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-semibold">Go-Insight</span>
            <span className="truncate text-xs text-muted-foreground">
              Observability Platform
            </span>
          </div>
        </div>
      </SidebarHeader>
      
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navigationItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton 
                    asChild 
                    isActive={pathname === item.url}
                  >
                    <Link href={item.url} className="flex items-center gap-2">
                      <item.icon className="h-4 w-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t border-sidebar-border">
        <div className="flex items-center justify-between px-4 py-2">
          <div className="flex items-center gap-2">
            <Settings className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">Settings</span>
          </div>
          <ThemeToggle />
        </div>
      </SidebarFooter>
      
      <SidebarRail />
    </Sidebar>
  )
}

interface DashboardLayoutProps {
  children: React.ReactNode
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <AppSidebar />
        <main className="flex-1 overflow-auto">
          <div className="container mx-auto p-6">
            {children}
          </div>
        </main>
      </div>
    </SidebarProvider>
  )
}