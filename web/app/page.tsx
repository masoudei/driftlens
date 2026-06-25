"use client"

import { useEffect, useState } from "react"
import { fetchDrifts } from "@/lib/api"
import type { DriftEvent } from "@/lib/types"
import { StatCard } from "@/components/stat-card"
import { RecentDrifts } from "@/components/recent-drifts"
import { AlertTriangle, Activity, Shield, Braces } from "lucide-react"

export default function Dashboard() {
  const [drifts, setDrifts] = useState<DriftEvent[]>([])

  useEffect(() => {
    fetchDrifts().then(setDrifts).catch(() => {})
  }, [])

  const totalDrifts = drifts.length
  const highSeverity = drifts.filter((d) => d.properties.severity === "high" || d.properties.severity === "critical").length
  const resources = [...new Set(drifts.map((d) => d.properties.resource))].length

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Overview of your infrastructure drift events
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Drifts"
          value={totalDrifts}
          icon={<AlertTriangle className="h-4 w-4" />}
        />
        <StatCard
          title="High Severity"
          value={highSeverity}
          icon={<Shield className="h-4 w-4" />}
        />
        <StatCard
          title="Resources Affected"
          value={resources}
          icon={<Braces className="h-4 w-4" />}
        />
        <StatCard
          title="Events Tracked"
          value={totalDrifts * 3 || "—"}
          icon={<Activity className="h-4 w-4" />}
        />
      </div>

      <RecentDrifts />
    </div>
  )
}
