import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { ScrollArea } from "@/components/ui/scroll-area";
import { SidebarGroupContent, SidebarMenu, SidebarMenuAction, SidebarMenuButton, SidebarMenuItem, SidebarMenuSub, SidebarMenuSubItem } from "@/components/ui/sidebar";
import { useYearGetAll } from "@/lib/api/year";
import { Link } from "@tanstack/react-router";
import { Calendar, ChevronRight } from "lucide-react";

function NavEvents() {
  const { data: years } = useYearGetAll();

  return (
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Events">
              <Link to="/events" activeProps={{ className: "border-b-2 border-b-primary" }}>
                <Calendar />
                <span>Events</span>
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
                            <Link to="/events/$year" params={{ year: year.formatted }} className="rounded-none border-l-2" activeProps={{ className: "border-l-primary" }}>
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

export default NavEvents;
