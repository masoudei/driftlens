"use client"

import { useSidebar } from "@/components/sidebar-context"
import { cn } from "@/lib/utils"
import type { ReactNode } from "react"

export function MainContent({ children }: { children: ReactNode }) {
  const { collapsed } = useSidebar()

  return (
    <div
      className={cn(
        "transition-all duration-300",
        "md:pl-[240px]",
        collapsed && "md:pl-[60px]"
      )}
    >
      {children}
    </div>
  )
}
