"use client"

import { useEffect, useState } from "react"
import { Badge } from "@/components/ui/badge"
import { buttonVariants } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { fetchDrifts } from "@/lib/api"
import type { DriftEvent } from "@/lib/types"
import { AlertTriangle, ArrowRight } from "lucide-react"
import Link from "next/link"

const severityColor: Record<string, string> = {
  low: "bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20",
  medium: "bg-amber-500/10 text-amber-500 hover:bg-amber-500/20",
  high: "bg-red-500/10 text-red-500 hover:bg-red-500/20",
  critical: "bg-red-600/10 text-red-600 hover:bg-red-600/20",
}

export function RecentDrifts() {
  const [drifts, setDrifts] = useState<DriftEvent[]>([])

  useEffect(() => {
    fetchDrifts().then(setDrifts).catch(() => {})
  }, [])

  const recent = drifts.slice(0, 8)

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>Latest Drifts</CardTitle>
          <CardDescription>Recent drift events across all resources</CardDescription>
        </div>
        <Link href="/drifts" className={buttonVariants({ variant: "ghost", size: "sm" })}>
          View all
          <ArrowRight className="ml-1 h-4 w-4" />
        </Link>
      </CardHeader>
      <CardContent>
        {recent.length === 0 ? (
          <div className="flex flex-col items-center gap-2 py-8 text-muted-foreground">
            <AlertTriangle className="h-8 w-8" />
            <p className="text-sm">No drifts detected yet</p>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Resource</TableHead>
                <TableHead>Name</TableHead>
                <TableHead>Drift</TableHead>
                <TableHead>Severity</TableHead>
                <TableHead>Detected</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {recent.map((d) => (
                <TableRow key={d.id}>
                  <TableCell className="font-medium">{d.properties.resource}</TableCell>
                  <TableCell>{d.properties.name}</TableCell>
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
                  <TableCell className="text-muted-foreground text-xs">
                    {new Date(d.properties.detected).toLocaleString()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
