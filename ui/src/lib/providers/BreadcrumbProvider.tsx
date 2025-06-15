import type { ReactNode } from "react";
import type { BreadcrumbAction, BreadcrumbState } from "../contexts/breadcrumbContext";
import { useMemo, useReducer } from "react";
import { BreadcrumbContext } from "../contexts/breadcrumbContext";

function breadcrumbReducer(state: BreadcrumbState, action: BreadcrumbAction): BreadcrumbState {
  switch (action.type) {
    case "ADD":
      if (state.some(crumb => crumb.title === action.payload.title)) {
        return state;
      }

      return [...state, action.payload];

    case "REMOVE":
      return state.filter(crumb => crumb.title !== action.payload.title);

    default:
      return state;
  }
}

export function BreadcrumbProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(breadcrumbReducer, []);

  const value = useMemo(() => ({ state, dispatch }), [state, dispatch]);

  return (
    <BreadcrumbContext value={value}>
      {children}
    </BreadcrumbContext>
  );
}
