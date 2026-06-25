"use client"

import { useEffect, useState } from "react"
import { fetchArgoApplications } from "@/lib/api"
import type { ArgoApp } from "@/lib/types"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { GitFork, Activity, AlertTriangle, ShieldCheck, RefreshCw } from "lucide-react"
import Link from "next/link"

const healthColors: Record<string, string> = {
  Healthy: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
  Degraded: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  Progressing: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
  Missing: "bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400",
  Suspended: "bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400",
  Unknown: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
}

const syncColors: Record<string, string> = {
  Synced: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
  OutOfSync: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
  Syncing: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
}

export default function ArgoCDPage() {
  const [apps, setApps] = useState<ArgoApp[]>([])
  const [loading, setLoading] = useState(true)

  function loadApps() {
    setLoading(true)
    fetchArgoApplications()
      .then(setApps)
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadApps() }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24 text-muted-foreground">
        <RefreshCw className="h-5 w-5 mr-2 animate-spin" />
        Loading ArgoCD applications...
      </div>
    )
  }

  if (apps.length === 0) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">ArgoCD</h1>
            <p className="text-sm text-muted-foreground mt-1">Tracked ArgoCD applications</p>
          </div>
        </div>
        <div className="flex flex-col items-center gap-2 py-16 text-muted-foreground">
          <GitFork className="h-8 w-8" />
          <p className="text-sm">No ArgoCD applications tracked yet</p>
          <p className="text-xs text-muted-foreground">
            Start DriftLens with DRIFTLENS_ARGOCD_URL set to connect to ArgoCD
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">ArgoCD</h1>
          <p className="text-sm text-muted-foreground mt-1">
            {apps.length} application{apps.length !== 1 ? "s" : ""} tracked
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadApps}>
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {apps.map((app) => (
          <Card key={app.app} className="hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <GitFork className="h-4 w-4 text-purple-500" />
                  <CardTitle className="text-base">{app.app}</CardTitle>
                </div>
                <Badge className={healthColors[app.health] || healthColors.Unknown}>
                  {app.health || "Unknown"}
                </Badge>
              </div>
              <CardDescription>
                <Badge variant="secondary" className={syncColors[app.sync_status]}>
                  {app.sync_status || "Unknown"}
                </Badge>
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between text-xs text-muted-foreground">
                <div className="flex items-center gap-1">
                  <AlertTriangle className="h-3 w-3" />
                  <span>{app.drift_count} drifts</span>
                </div>
                <div className="flex items-center gap-1">
                  <Activity className="h-3 w-3" />
                  <span>{app.action || "-"}</span>
                </div>
                <div className="flex items-center gap-1">
                  <ShieldCheck className="h-3 w-3" />
                  <span>Risk: {app.drift_count > 5 ? "High" : app.drift_count > 2 ? "Medium" : "Low"}</span>
                </div>
              </div>
              <div className="mt-3 flex gap-2">
                <Link href={`/timeline/${app.app}`} className="flex-1">
                  <Button variant="outline" size="sm" className="w-full text-xs">Timeline</Button>
                </Link>
                <Link href={`/risk/${app.app}`} className="flex-1">
                  <Button variant="outline" size="sm" className="w-full text-xs">Risk</Button>
                </Link>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {apps.length > 0 && (
        <p className="text-xs text-muted-foreground text-center">
          Last updated: {new Date().toLocaleString()}
        </p>
      )}
    </div>
  )
}
