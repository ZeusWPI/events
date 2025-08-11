import { Link, useMatch } from "@tanstack/react-router";
import { CalendarIcon } from "lucide-react";
import { Collapsible } from "@/components/ui/collapsible";
import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";

function NavEvents() {
  const isActive = useMatch({ from: "/events", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Events">
              <Link to="/events">
                <CalendarIcon className={isActive && "stroke-primary"} />
                <span>Events</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </Collapsible>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}

export default NavEvents;
