import EventsIcon from "@/components/icons/EventsIcon";
import { YearSelector } from "@/components/molecules/YearSelector";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton, SidebarMenuItem, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { useBreadcrumbs } from "@/lib/hooks/useBreadcrumb";
import { formatDate, getBuildTime } from "@/lib/utils/utils";
import { Link } from "@tanstack/react-router";
import type { ReactNode } from "react";
import { Fragment } from "react";
import { NavAnnouncements } from "./NavAnnouncements";
import NavEvents from "./NavEvents";
import { NavMails } from "./NavMails";
import { NavPowerpoints } from "./NavPowerpoints";
import NavTasks from "./NavTasks";
import { NavUser } from "./NavUser";
import { NavDeadlines } from "./NavDeadlines";

function AppSidebar({ children }: { children: ReactNode }) {
  const { state: breadcrumbs } = useBreadcrumbs();

  const buildTime = getBuildTime()

  return (
    <SidebarProvider>
      <Sidebar variant="inset" collapsible="icon">
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton size="lg" asChild>
                <Link to="/">
                  <div className="flex aspect-square size-8 items-center justify-center ">
                    <EventsIcon className="[&:not([data-state=open])]:ml-0.5 size-8 fill-secondary-foreground" />
                  </div>
                  <div className="flex flex-col gap-0.5 leading-none">
                    <span className="font-semibold">Events</span>
                    <span className="text-xs text-muted-foreground">
                      {`Built: ${buildTime ? formatDate(buildTime) : "Development"}`}
                    </span>
                  </div>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent className="pt-[16px]">
          <SidebarGroup>
            <NavEvents />
            <NavAnnouncements />
            <NavMails />
            <NavPowerpoints />
            <NavDeadlines />
            <NavTasks />
          </SidebarGroup>
          <SidebarGroup className="mt-auto">
            <YearSelector />
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <NavUser />
        </SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <div className="container mx-auto p-2 h-full">
          <header className="flex h-16 shrink-0 items-center gap-2">
            <SidebarTrigger />
            <Separator orientation="vertical" className="mr-2 h-4" />
            {/* FIX: Can somtimes span multiple lines */}
            <Breadcrumb>
              <BreadcrumbList>
                {breadcrumbs.map((crumb, index) => (
                  <Fragment key={crumb.title}>
                    <BreadcrumbItem>
                      {crumb.link
                        ? (
                          <BreadcrumbLink asChild>
                            <Link to={crumb.link.to} params={crumb.link.params}>
                              {crumb.title}
                            </Link>
                          </BreadcrumbLink>
                        )
                        : crumb.title}
                    </BreadcrumbItem>
                    {index !== breadcrumbs.length - 1 && <BreadcrumbSeparator />}
                  </Fragment>
                ))}
              </BreadcrumbList>
            </Breadcrumb>
          </header>

          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

export default AppSidebar;
