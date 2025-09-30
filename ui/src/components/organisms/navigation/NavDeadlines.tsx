import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { Link, useMatch } from "@tanstack/react-router";
import { ClockAlertIcon } from "lucide-react";

import { Collapsible } from "@/components/ui/collapsible";

export function NavDeadlines() {
  const isActive = useMatch({ from: "/deadlines", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Deadlines">
              <Link to="/deadlines">
                <ClockAlertIcon className={isActive && "stroke-primary"} />
                <span>Deadlines</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </Collapsible>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}
