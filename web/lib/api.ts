import type { DriftEvent, TimelineNode, RiskResponse } from "./types"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "/api"

async function fetchJSON<T>(url: string): Promise<T> {
  const res = await fetch(url)
  if (!res.ok) {
    throw new Error(`API ${res.status}: ${res.statusText}`)
  }
  return res.json()
}

export async function fetchDrifts(): Promise<DriftEvent[]> {
  const data = await fetchJSON<{ drifts: DriftEvent[] }>(`${API_BASE}/drifts`)
  return data.drifts
}

export async function fetchDrift(id: string): Promise<DriftEvent | null> {
  try {
    const data = await fetchJSON<{ drift: DriftEvent }>(`${API_BASE}/drifts/${id}`)
    return data.drift
  } catch {
    return null
  }
}

export async function fetchTimeline(resource: string): Promise<TimelineNode[]> {
  const data = await fetchJSON<{ timeline: TimelineNode[] }>(
    `${API_BASE}/timeline/${resource}`
  )
  return data.timeline
}

export async function fetchRisk(resource: string): Promise<RiskResponse | null> {
  try {
    return await fetchJSON<RiskResponse>(`${API_BASE}/risk/${resource}`)
  } catch {
    return null
  }
}
