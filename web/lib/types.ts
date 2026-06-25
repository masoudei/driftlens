export interface DriftEvent {
  id: string
  properties: {
    resource: string
    name: string
    namespace: string
    drift: string
    severity: string
    detected: string
  }
  sources: string[]
  timestamp: string
}

export interface TimelineNode {
  id: string
  type: string
  props: Record<string, string>
  timestamp: string
  chain: string | null
}

export interface TimelineResponse {
  timeline: TimelineNode[]
}

export interface RiskResponse {
  resource: string
  drift_count: number
  risk_severity: string
}

export interface DriftsResponse {
  drifts: DriftEvent[]
}

export interface HealthResponse {
  status: string
}

export type ArgoApp = {
  app: string
  sync_status: string
  health: string
  action: string
  last_seen: string
  drift_count: number
}

export type NavItem = {
  title: string
  href: string
  icon: string
}
