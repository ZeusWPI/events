import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { } from "@/lib/api/year";
import { Link, useMatch } from "@tanstack/react-router";
import { MegaphoneIcon } from "lucide-react";

export function NavAnnouncements() {
  const isActive = useMatch({ from: "/announcements", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton asChild tooltip="Announcements">
            <Link to="/announcements">
              <MegaphoneIcon className={isActive && "stroke-primary"} />
              <span>Announcements</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}
