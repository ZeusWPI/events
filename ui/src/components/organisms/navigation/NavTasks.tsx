import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { Link, useMatch } from "@tanstack/react-router";
import { ClipboardCheck } from "lucide-react";

function NavTasks() {
  const isActive = useMatch({ from: "/tasks", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton asChild tooltip="Tasks">
            <Link to="/tasks">
              <ClipboardCheck className={isActive && "stroke-primary"} />
              <span>Tasks</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}

export default NavTasks;
