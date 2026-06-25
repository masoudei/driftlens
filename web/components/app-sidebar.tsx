"use client"

import { cn } from "@/lib/utils"
import {
  LayoutDashboard,
  AlertTriangle,
  Timeline,
  Shield,
  Braces,
  PanelLeftClose,
  PanelLeft,
} from "lucide-react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet"
import { useSidebar } from "@/components/sidebar-context"

const navItems = [
  { title: "Dashboard", href: "/", icon: "LayoutDashboard" },
  { title: "Drifts", href: "/drifts", icon: "AlertTriangle" },
  { title: "Timeline", href: "/timeline", icon: "Timeline" },
  { title: "Risk", href: "/risk", icon: "Shield" },
]

const resourceTypes = [
  { title: "Deployment", href: "/timeline/Deployment", icon: "Braces" },
  { title: "StatefulSet", href: "/timeline/StatefulSet", icon: "Braces" },
  { title: "DaemonSet", href: "/timeline/DaemonSet", icon: "Braces" },
  { title: "ConfigMap", href: "/timeline/ConfigMap", icon: "Braces" },
  { title: "Secret", href: "/timeline/Secret", icon: "Braces" },
  { title: "Namespace", href: "/timeline/Namespace", icon: "Braces" },
]

const iconMap: Record<string, React.ElementType> = {
  LayoutDashboard,
  AlertTriangle,
  Timeline,
  Shield,
  Braces,
}

function NavContent({ collapsed, onNavigate }: { collapsed: boolean; onNavigate?: () => void }) {
  const pathname = usePathname()

  return (
    <nav className="flex flex-col gap-1">
      {!collapsed && (
        <span className="px-2 py-1 text-xs font-medium text-muted-foreground">Overview</span>
      )}
      {navItems.map((item) => {
        const Icon = iconMap[item.icon]
        const active = item.href === "/" ? pathname === "/" : pathname.startsWith(item.href)
        return (
          <Link
            key={item.href}
            href={item.href}
            onClick={onNavigate}
            className={cn(
              "flex h-9 items-center gap-3 rounded-md px-2 text-sm font-medium transition-colors",
              active
                ? "bg-sidebar-accent text-sidebar-accent-foreground"
                : "text-sidebar-foreground hover:bg-sidebar-accent/50"
            )}
          >
            {Icon && <Icon className="h-4 w-4 shrink-0" />}
            {!collapsed && item.title}
          </Link>
        )
      })}

      {!collapsed && (
        <>
          <Separator className="my-3" />
          <span className="px-2 py-1 text-xs font-medium text-muted-foreground">Resources</span>
        </>
      )}
      {collapsed && <Separator className="my-3" />}

      {resourceTypes.map((item) => {
        const Icon = iconMap[item.icon]
        const active = pathname === item.href
        return (
          <Link
            key={item.href}
            href={item.href}
            onClick={onNavigate}
            className={cn(
              "flex h-9 items-center gap-3 rounded-md px-2 text-sm font-medium transition-colors",
              active
                ? "bg-sidebar-accent text-sidebar-accent-foreground"
                : "text-sidebar-foreground hover:bg-sidebar-accent/50"
            )}
          >
            {Icon && <Icon className="h-4 w-4 shrink-0" />}
            {!collapsed && item.title}
          </Link>
        )
      })}
    </nav>
  )
}

export function AppSidebar() {
  const { collapsed, setCollapsed, mobileOpen, setMobileOpen } = useSidebar()

  return (
    <>
      <aside
        className={cn(
          "fixed left-0 top-0 z-30 hidden h-full flex-col border-r bg-sidebar transition-all duration-300 md:flex",
          collapsed ? "w-[60px]" : "w-[240px]"
        )}
      >
        <div className="flex h-14 items-center gap-2 border-b px-4">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground text-xs font-bold">
            DL
          </div>
          {!collapsed && (
            <span className="font-semibold text-sm tracking-tight">DriftLens</span>
          )}
        </div>

        <ScrollArea className="flex-1 px-3 py-3">
          <NavContent collapsed={collapsed} />
        </ScrollArea>

        <div className="border-t p-3">
          <Button
            variant="ghost"
            size="sm"
            className="w-full justify-start text-muted-foreground hover:text-foreground"
            onClick={() => setCollapsed(!collapsed)}
          >
            {collapsed ? (
              <PanelLeft className="h-4 w-4" />
            ) : (
              <>
                <PanelLeftClose className="h-4 w-4 mr-2" />
                Collapse
              </>
            )}
          </Button>
        </div>
      </aside>

      <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
        <SheetContent side="left" className="w-[240px] p-0">
          <SheetHeader className="flex h-14 flex-row items-center gap-2 border-b px-4">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground text-xs font-bold">
              DL
            </div>
            <SheetTitle className="font-semibold text-sm tracking-tight">DriftLens</SheetTitle>
          </SheetHeader>
          <ScrollArea className="flex-1 px-3 py-3 h-[calc(100vh-3.5rem)]">
            <NavContent collapsed={false} onNavigate={() => setMobileOpen(false)} />
          </ScrollArea>
        </SheetContent>
      </Sheet>
    </>
  )
}
