"use client"

import { useParams } from "next/navigation"
import { TimelineView } from "@/components/timeline-view"
import { buttonVariants } from "@/components/ui/button"
import { ArrowLeft } from "lucide-react"
import Link from "next/link"

export default function ResourceTimelinePage() {
  const params = useParams()
  const resource = params.resource as string

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/timeline" className={buttonVariants({ variant: "ghost", size: "icon" })}>
          <ArrowLeft className="h-4 w-4" />
        </Link>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">{resource}</h1>
          <p className="text-sm text-muted-foreground mt-1">Event timeline</p>
        </div>
      </div>
      <TimelineView resource={resource} />
    </div>
  )
}
