import { SidebarGroupContent, SidebarMenu, SidebarMenuItem, SidebarMenuButton } from "@/components/ui/sidebar";
import { Link, useMatch } from "@tanstack/react-router";
import { MonitorStopIcon } from "lucide-react";

import { Collapsible } from "@/components/ui/collapsible";

export function NavPowerpoints() {
  const isActive = useMatch({ from: "/powerpoints", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Powerpoints">
              <Link to="/powerpoints">
                <MonitorStopIcon className={isActive && "stroke-primary"} />
                <span>Powerpoints</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </Collapsible>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}
