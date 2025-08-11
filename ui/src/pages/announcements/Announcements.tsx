import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { weightCategory } from "@/lib/types/general";
import { Outlet } from "@tanstack/react-router";

export function Announcements() {
  useBreadcrumb({ title: "Announcements", weight: weightCategory, link: { to: "/announcements" } });

  return <Outlet />
}
