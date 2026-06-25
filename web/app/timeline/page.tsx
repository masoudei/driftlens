"use client"

import { useState } from "react"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { TimelineView } from "@/components/timeline-view"

const resourceTypes = ["Deployment", "StatefulSet", "DaemonSet", "ConfigMap", "Secret", "Namespace", "ArgoSync"]

export default function TimelinePage() {
  const [resource, setResource] = useState("Deployment")

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Timeline</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Event timeline for a specific resource type
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

      <TimelineView resource={resource} />
    </div>
  )
}
