import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { weightCategory } from "@/lib/types/general";
import { Outlet } from "@tanstack/react-router";

function Events() {
  useBreadcrumb({ title: "Events", weight: weightCategory, link: { to: "/events" } });

  return <Outlet />
}

export default Events;
