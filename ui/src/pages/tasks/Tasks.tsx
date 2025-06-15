import { Outlet } from "@tanstack/react-router";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";

export function Tasks() {
  useBreadcrumb({ title: "Tasks", link: { to: "/tasks" } });

  return (
    <Outlet />
  );
}
