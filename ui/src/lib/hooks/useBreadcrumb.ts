import type { Breadcrumb } from "../context/breadcrumbContext"
import { useContext, useEffect } from "react"
import { BreadcrumbContext } from "../context/breadcrumbContext"

export function useBreadcrumb(crumb: Breadcrumb) {
  const context = useContext(BreadcrumbContext)
  if (!context) {
    throw new Error("useBreadcrumbs must be used within a BreadcrumbProvider")
  }

  useEffect(() => {
    context.dispatch({ type: "ADD", payload: crumb })

    return () => {
      context.dispatch({ type: "REMOVE", payload: crumb })
    }

    // Use stringify to avoid infinite rerenders
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(crumb), context.dispatch])
}
