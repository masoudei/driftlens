"use client"

import { useState } from "react"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { RiskView } from "@/components/risk-view"

const resourceTypes = ["Deployment", "StatefulSet", "DaemonSet", "ConfigMap", "Secret", "Namespace"]

export default function RiskPage() {
  const [resource, setResource] = useState("Deployment")

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Risk Analysis</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Risk assessment for each resource type
          </p>
        </div>
        <Select value={resource} onValueChange={(v) => v && setResource(v)}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Resource type" />
          </SelectTrigger>
          <SelectContent>
            {resourceTypes.map((r) => (
              <SelectItem key={r} value={r}>{r}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <RiskView resource={resource} />
    </div>
  )
}
