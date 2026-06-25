"use client"

import { useParams } from "next/navigation"
import Link from "next/link"
import { ArrowLeft } from "lucide-react"
import { buttonVariants } from "@/components/ui/button"
import { TimelineView } from "@/components/timeline-view"

export default function ArgoAppPage() {
  const params = useParams()
  const app = params.app as string

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/argocd" className={buttonVariants({ variant: "ghost", size: "icon" })}>
          <ArrowLeft className="h-4 w-4" />
        </Link>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">{app}</h1>
          <p className="text-sm text-muted-foreground mt-1">ArgoCD application timeline</p>
        </div>
      </div>
      <TimelineView resource={app} />
    </div>
  )
}
