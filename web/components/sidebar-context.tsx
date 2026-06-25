"use client"

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from "react"

interface SidebarContextType {
  collapsed: boolean
  setCollapsed: (v: boolean) => void
  mobileOpen: boolean
  setMobileOpen: (v: boolean) => void
}

const SidebarContext = createContext<SidebarContextType>({
  collapsed: false,
  setCollapsed: () => {},
  mobileOpen: false,
  setMobileOpen: () => {},
})

export function SidebarProvider({ children }: { children: ReactNode }) {
  const [collapsed, setCollapsed] = useState(() => {
    if (typeof window === "undefined") return false
    const w = window.innerWidth
    return w >= 768 && w < 1024
  })
  const [mobileOpen, setMobileOpen] = useState(false)

  const handleResize = useCallback(() => {
    const w = window.innerWidth
    if (w < 768) {
      setCollapsed(false)
    } else if (w < 1024) {
      setCollapsed(true)
    } else {
      setCollapsed(false)
    }
  }, [])

  useEffect(() => {
    window.addEventListener("resize", handleResize)
    return () => window.removeEventListener("resize", handleResize)
  }, [handleResize])

  return (
    <SidebarContext.Provider value={{ collapsed, setCollapsed, mobileOpen, setMobileOpen }}>
      {children}
    </SidebarContext.Provider>
  )
}

export function useSidebar() {
  return useContext(SidebarContext)
}
