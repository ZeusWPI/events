import { Link, useMatch } from "@tanstack/react-router";
import { ChevronRight, MegaphoneIcon } from "lucide-react";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { ScrollArea } from "@/components/ui/scroll-area";
import { SidebarGroupContent, SidebarMenu, SidebarMenuAction, SidebarMenuButton, SidebarMenuItem, SidebarMenuSub, SidebarMenuSubItem } from "@/components/ui/sidebar";
import { useYearGetAll } from "@/lib/api/year";

export function NavAnnouncements() {
  const isActive = useMatch({ from: "/announcements", shouldThrow: false });

  const { data: years } = useYearGetAll();

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Announcements">
              <Link to="/announcements">
                <MegaphoneIcon className={isActive && "stroke-primary"} />
                <span>Announcements</span>
              </Link>
            </SidebarMenuButton>
            {years && (
              <>
                <CollapsibleTrigger asChild>
                  <SidebarMenuAction className="data-[state=open]:rotate-90">
                    <ChevronRight />
                  </SidebarMenuAction>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub className="pl-0 py-0 gap-0 border-none">
                    <ScrollArea className="h-44">
                      {years.map(year => (
                        <SidebarMenuSubItem key={year.id}>
                          <SidebarMenuButton asChild>
                            <Link to="/announcements/$year" params={{ year: year.formatted }} className="rounded-none border-l-2" activeProps={{ className: "border-l-primary" }}>
                              <span>{year.formatted}</span>
                            </Link>
                          </SidebarMenuButton>
                        </SidebarMenuSubItem>
                      ))}
                    </ScrollArea>
                  </SidebarMenuSub>
                </CollapsibleContent>
              </>
            )}
          </SidebarMenuItem>
        </Collapsible>
      </SidebarMenu>
    </SidebarGroupContent>
  );
}
