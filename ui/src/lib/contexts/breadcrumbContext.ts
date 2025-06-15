import type { LinkProps } from "@tanstack/react-router";
import { createContext } from "react";

export interface Breadcrumb {
  title: string;
  link?: LinkProps;
}

export type BreadcrumbAction
  = | { type: "ADD"; payload: Breadcrumb }
    | { type: "REMOVE"; payload: Breadcrumb };

export type BreadcrumbState = Breadcrumb[];

export const BreadcrumbContext = createContext<{
  state: BreadcrumbState;
  dispatch: React.Dispatch<BreadcrumbAction>;
} | undefined>(undefined);
