import { Link, useMatch } from "@tanstack/react-router";
import { MailIcon } from "lucide-react";
import { Collapsible } from "@/components/ui/collapsible";
import { SidebarGroupContent, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";

export function NavMails() {
  const isActive = useMatch({ from: "/mails", shouldThrow: false });

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Mails">
              <Link to="/mails">
                <MailIcon className={isActive && "stroke-primary"} />
                <span>Mails</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </Collapsible>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}
