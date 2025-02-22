import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { Link } from "@tanstack/react-router";
import { CalendarDays } from "lucide-react";

function NavEvents() {
  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton asChild tooltip="Overview">
            <Link to="/events" mask={{ to: "/" }} activeProps={{ className: "border-b-2 border-b-primary" }} className="">
              <CalendarDays />
              <span>Overview</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}

export default NavEvents;
