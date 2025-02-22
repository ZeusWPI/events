import type { Breadcrumb } from "../context/breadcrumbContext";
import { useContext, useEffect } from "react";
import { BreadcrumbContext } from "../context/breadcrumbContext";

export function useBreadcrumbs() {
  const context = useContext(BreadcrumbContext);
  if (!context) {
    throw new Error("useBreadcrumbs must be used within a BreadcrumbProvider");
  }
  return context;
}

export function useBreadcrumb(crumb: Breadcrumb) {
  const { dispatch } = useBreadcrumbs();

  useEffect(() => {
    dispatch({ type: "ADD", payload: crumb });

    return () => {
      dispatch({ type: "REMOVE", payload: crumb });
    };

    // Use stringify to avoid infinite rerenders
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(crumb), dispatch]);
}
