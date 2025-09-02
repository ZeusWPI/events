import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { Link, useMatch } from "@tanstack/react-router";
import { CalendarIcon } from "lucide-react";

function NavEvents() {
  const isActive = useMatch({ from: "/events", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton asChild tooltip="Events">
            <Link to="/events">
              <CalendarIcon className={isActive && "stroke-primary"} />
              <span>Events</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}

export default NavEvents;
