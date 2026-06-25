"use client"

import { useEffect, useState } from "react"
import { fetchTimeline } from "@/lib/api"
import type { TimelineNode } from "@/lib/types"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { AlertTriangle, Activity, Circle } from "lucide-react"

interface TimelineViewProps {
  resource: string
}

const typeConfig: Record<string, { label: string; color: string; icon: React.ElementType }> = {
  KubernetesEvent: {
    label: "K8s Event",
    color: "border-l-blue-500",
    icon: Activity,
  },
  DriftEvent: {
    label: "Drift",
    color: "border-l-amber-500",
    icon: AlertTriangle,
  },
}

export function TimelineView({ resource }: TimelineViewProps) {
  const [events, setEvents] = useState<TimelineNode[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchTimeline(resource)
      .then((data) => {
        setEvents(data.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()))
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [resource])

  if (loading) {
    return <div className="text-muted-foreground text-sm py-8 text-center">Loading timeline...</div>
  }

  if (events.length === 0) {
    return (
      <div className="flex flex-col items-center gap-2 py-16 text-muted-foreground">
        <Activity className="h-8 w-8" />
        <p className="text-sm">No events found for {resource}</p>
      </div>
    )
  }

  const groups: Record<string, string> = {}
  events.forEach((e) => {
    if (typeof e.chain === "string") {
      const [chainId] = e.chain.split(":")
      groups[e.id] = chainId
    }
  })

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm font-medium">Timeline</CardTitle>
        <CardDescription>
          {events.length} events for {resource}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="relative space-y-0">
          {events.map((event, i) => {
            const config = typeConfig[event.type] || {
              label: event.type,
              color: "border-l-gray-500",
              icon: Circle,
            }
            const Icon = config.icon
            const isLast = i === events.length - 1

            return (
              <div key={event.id} className="relative flex gap-4 pb-6">
                <div className="flex flex-col items-center">
                  <div
                    className={`flex h-8 w-8 items-center justify-center rounded-full border bg-background ${config.color.replace("border-l-", "border-")}`}
                  >
                    <Icon className="h-4 w-4" />
                  </div>
                  {!isLast && (
                    <div className="mt-1 h-full w-px bg-border" />
                  )}
                </div>
                <div className="flex-1 min-w-0 pt-1">
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className="text-xs">
                      {config.label}
                    </Badge>
                    <span className="text-xs text-muted-foreground">
                      {new Date(event.timestamp).toLocaleString()}
                    </span>
                  </div>
                  <div className="mt-1 space-y-0.5">
                    {Object.entries(event.props).map(([k, v]) => (
                      <p key={k} className="text-xs text-muted-foreground">
                        <span className="font-medium text-foreground">{k}:</span>{" "}
                        <span className="font-mono">{v}</span>
                      </p>
                    ))}
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}
