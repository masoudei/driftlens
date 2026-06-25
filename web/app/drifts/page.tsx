"use client"

import { useEffect, useState } from "react"
import { fetchDrifts } from "@/lib/api"
import type { DriftEvent } from "@/lib/types"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Card,
  CardContent,
  CardHeader,
} from "@/components/ui/card"
import { Search } from "lucide-react"
import Link from "next/link"
import { buttonVariants } from "@/components/ui/button"

const severityColor: Record<string, string> = {
  low: "bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20",
  medium: "bg-amber-500/10 text-amber-500 hover:bg-amber-500/20",
  high: "bg-red-500/10 text-red-500 hover:bg-red-500/20",
  critical: "bg-red-600/10 text-red-600 hover:bg-red-600/20",
}

export default function DriftsPage() {
  const [drifts, setDrifts] = useState<DriftEvent[]>([])
  const [search, setSearch] = useState("")
  const [severityFilter, setSeverityFilter] = useState("all")
  const [resourceFilter, setResourceFilter] = useState("all")

  useEffect(() => {
    fetchDrifts().then(setDrifts).catch(() => {})
  }, [])

  const resources = [...new Set(drifts.map((d) => d.properties.resource))]

  const filtered = drifts.filter((d) => {
    if (severityFilter !== "all" && d.properties.severity !== severityFilter) return false
    if (resourceFilter !== "all" && d.properties.resource !== resourceFilter) return false
    if (search) {
      const q = search.toLowerCase()
      return (
        d.properties.name.toLowerCase().includes(q) ||
        d.properties.resource.toLowerCase().includes(q) ||
        d.properties.namespace.toLowerCase().includes(q) ||
        d.properties.drift.toLowerCase().includes(q)
      )
    }
    return true
  })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Drifts</h1>
        <p className="text-sm text-muted-foreground mt-1">
          All detected drift events across your infrastructure
        </p>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap items-center gap-2">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search drifts..."
                  className="w-[200px] sm:w-[250px] pl-9 h-9"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>
              <Select value={severityFilter} onValueChange={(v) => v && setSeverityFilter(v)}>
                <SelectTrigger className="w-[130px] h-9">
                  <SelectValue placeholder="Severity" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Severities</SelectItem>
                  <SelectItem value="low">Low</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                  <SelectItem value="critical">Critical</SelectItem>
                </SelectContent>
              </Select>
              <Select value={resourceFilter} onValueChange={(v) => v && setResourceFilter(v)}>
                <SelectTrigger className="w-[150px] h-9">
                  <SelectValue placeholder="Resource" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Resources</SelectItem>
                  {resources.map((r) => (
                    <SelectItem key={r} value={r}>{r}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <p className="text-xs text-muted-foreground">
              {filtered.length} of {drifts.length} drifts
            </p>
          </div>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Resource</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Namespace</TableHead>
                  <TableHead>Drift</TableHead>
                  <TableHead>Severity</TableHead>
                  <TableHead>Detected</TableHead>
                  <TableHead></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-muted-foreground py-8">
                      No drifts found
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((d) => (
                    <TableRow key={d.id}>
                      <TableCell className="font-medium">{d.properties.resource}</TableCell>
                      <TableCell>{d.properties.name}</TableCell>
                      <TableCell className="text-muted-foreground">{d.properties.namespace}</TableCell>
                      <TableCell className="max-w-[200px] truncate font-mono text-xs">
                        {d.properties.drift}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className={severityColor[d.properties.severity] || ""}
                        >
                          {d.properties.severity}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-xs text-muted-foreground">
                        {new Date(d.properties.detected).toLocaleString()}
                      </TableCell>
                      <TableCell>
                        <Link
                          href={`/drifts/${encodeURIComponent(d.id)}`}
                          className={buttonVariants({ variant: "ghost", size: "sm" })}
                        >
                          View
                        </Link>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
