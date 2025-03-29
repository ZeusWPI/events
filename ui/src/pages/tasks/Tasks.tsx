import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Outlet } from "@tanstack/react-router";

export function Tasks() {
  useBreadcrumb({ title: "Tasks", link: { to: "/tasks" } });

  return (
    <Outlet />
  );
}
