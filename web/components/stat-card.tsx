import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import { ReactNode } from "react"

interface StatCardProps {
  title: string
  value: string | number
  change?: string
  changeDirection?: "up" | "down"
  icon: ReactNode
}

export function StatCard({ title, value, change, changeDirection, icon }: StatCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
          {icon}
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {change && (
          <p
            className={cn(
              "text-xs mt-1",
              changeDirection === "up" ? "text-emerald-500" : "text-red-500"
            )}
          >
            {change}
          </p>
        )}
      </CardContent>
    </Card>
  )
}
