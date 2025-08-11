import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { weightCategory } from "@/lib/types/general";
import { Outlet } from "@tanstack/react-router";

export function Mails() {
  useBreadcrumb({ title: "Mails", weight: weightCategory, link: { to: "/mails" } });

  return <Outlet />
}
