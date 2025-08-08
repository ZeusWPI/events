import { Outlet } from "@tanstack/react-router";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { weightCategory } from "@/lib/types/general";

export function Tasks() {
  useBreadcrumb({ title: "Tasks", weight: weightCategory, link: { to: "/tasks" } });

  return (
    <Outlet />
  );
}
