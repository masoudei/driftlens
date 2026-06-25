"use client"

import { useEffect, useState } from "react"
import { fetchRisk, fetchDrifts } from "@/lib/api"
import type { RiskResponse, DriftEvent } from "@/lib/types"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Shield, ShieldAlert, ShieldCheck } from "lucide-react"

interface RiskViewProps {
  resource: string
}

const severityColor: Record<string, string> = {
  low: "bg-emerald-500/10 text-emerald-500",
  medium: "bg-amber-500/10 text-amber-500",
  high: "bg-red-500/10 text-red-500",
  critical: "bg-red-600/10 text-red-600",
}

const severityProgress: Record<string, number> = {
  low: 25,
  medium: 50,
  high: 75,
  critical: 100,
}

export function RiskView({ resource }: RiskViewProps) {
  const [risk, setRisk] = useState<RiskResponse | null>(null)
  const [drifts, setDrifts] = useState<DriftEvent[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    Promise.all([
      fetchRisk(resource),
      fetchDrifts(),
    ])
      .then(([r, d]) => {
        setRisk(r)
        setDrifts(d.filter((dr) => dr.properties.resource === resource))
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [resource])

  if (loading) {
    return <div className="text-muted-foreground text-sm py-8 text-center">Loading risk analysis...</div>
  }

  const severityCounts = {
    low: drifts.filter((d) => d.properties.severity === "low").length,
    medium: drifts.filter((d) => d.properties.severity === "medium").length,
    high: drifts.filter((d) => d.properties.severity === "high").length,
    critical: drifts.filter((d) => d.properties.severity === "critical").length,
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle className="text-sm font-medium">Overall Risk</CardTitle>
            <CardDescription>
              Risk assessment for {resource}
            </CardDescription>
          </div>
          {risk && (
            <Badge
              variant="secondary"
              className={severityColor[risk.risk_severity] || ""}
            >
              {risk.risk_severity.toUpperCase()}
            </Badge>
          )}
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center gap-4">
              {risk && (
                <>
                  {risk.risk_severity === "low" ? (
                    <ShieldCheck className="h-10 w-10 text-emerald-500" />
                  ) : risk.risk_severity === "medium" ? (
                    <Shield className="h-10 w-10 text-amber-500" />
                  ) : (
                    <ShieldAlert className="h-10 w-10 text-red-500" />
                  )}
                  <div className="flex-1">
                    <Progress
                      value={risk ? severityProgress[risk.risk_severity] || 0 : 0}
                      className="h-2"
                    />
                    <p className="text-xs text-muted-foreground mt-1">
                      {risk?.drift_count} drift event{risk?.drift_count !== 1 ? "s" : ""} detected
                    </p>
                  </div>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-4 md:grid-cols-4">
        {(["low", "medium", "high", "critical"] as const).map((severity) => (
          <Card key={severity}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium capitalize text-muted-foreground">
                {severity}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div
                className={`text-2xl font-bold ${
                  severity === "low"
                    ? "text-emerald-500"
                    : severity === "medium"
                    ? "text-amber-500"
                    : severity === "high"
                    ? "text-red-500"
                    : "text-red-600"
                }`}
              >
                {severityCounts[severity]}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
