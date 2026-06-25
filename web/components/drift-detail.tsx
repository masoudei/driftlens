"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import { fetchDrift, fetchTimeline } from "@/lib/api"
import type { DriftEvent, TimelineNode } from "@/lib/types"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ArrowLeft, ExternalLink } from "lucide-react"
import Link from "next/link"
import { buttonVariants } from "@/components/ui/button"

const severityColor: Record<string, string> = {
  low: "bg-emerald-500/10 text-emerald-500",
  medium: "bg-amber-500/10 text-amber-500",
  high: "bg-red-500/10 text-red-500",
  critical: "bg-red-600/10 text-red-600",
}

export function DriftDetail() {
  const params = useParams()
  const [drift, setDrift] = useState<DriftEvent | null>(null)
  const [timeline, setTimeline] = useState<TimelineNode[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const id = params.id as string
    Promise.all([
      fetchDrift(id),
      fetchTimeline("Deployment"),
    ]).then(([d, t]) => {
      setDrift(d)
      setTimeline(t)
      setLoading(false)
    })
  }, [params.id])

  if (loading) {
    return <div className="text-muted-foreground text-sm">Loading...</div>
  }

  if (!drift) {
    return (
      <div className="flex flex-col items-center gap-4 py-16">
        <p className="text-muted-foreground">Drift not found</p>
        <Link href="/drifts" className={buttonVariants({ variant: "outline" })}>
          Back to drifts
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/drifts" className={buttonVariants({ variant: "ghost", size: "icon" })}>
          <ArrowLeft className="h-4 w-4" />
        </Link>
        <div>
          <h2 className="text-lg font-semibold">Drift Details</h2>
          <p className="text-sm text-muted-foreground">{drift.id}</p>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Resource</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Type</span>
                <span className="font-medium">{drift.properties.resource}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Name</span>
                <span className="font-medium">{drift.properties.name}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Namespace</span>
                <span className="font-medium">{drift.properties.namespace}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Drift Info</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Severity</span>
                <Badge
                  variant="secondary"
                  className={severityColor[drift.properties.severity] || ""}
                >
                  {drift.properties.severity}
                </Badge>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Detected</span>
                <span className="font-medium">
                  {new Date(drift.properties.detected).toLocaleString()}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">Changes</CardTitle>
          <CardDescription>What properties changed</CardDescription>
        </CardHeader>
        <CardContent>
          <code className="block rounded-lg bg-muted p-4 text-sm font-mono">
            {drift.properties.drift}
          </code>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">Sources</CardTitle>
          <CardDescription>Related nodes in the causality graph</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {drift.sources.map((s) => (
              <div
                key={s}
                className="flex items-center gap-2 text-sm text-muted-foreground"
              >
                <ExternalLink className="h-3 w-3" />
                {s}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
