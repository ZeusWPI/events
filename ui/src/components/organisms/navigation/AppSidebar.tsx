import EventsIcon from "@/components/icons/EventsIcon";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton, SidebarMenuItem, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumbs } from "@/lib/hooks/useBreadcrumb";
import { useYear } from "@/lib/hooks/useYear";
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

function AppSidebar({ children }: { children: ReactNode }) {
  const { state: breadcrumbs } = useBreadcrumbs();

  const { year, setYear, locked } = useYear()
  const { data: years } = useYearGetAll()

  const handleSelectChange = (value: string) => {
    const newYear = years?.find(y => y.id === Number(value))
    if (!newYear || newYear?.id === year?.id) {
      return
    }

    setYear(newYear)
  }

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
                      {`Built: ${formatDate(getBuildTime())}`}
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
            <NavTasks />
          </SidebarGroup>
          <SidebarGroup className="mt-auto">
            <Select onValueChange={handleSelectChange} defaultValue={year?.id.toString()} disabled={locked}>
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent className="max-h-72">
                {years?.map(y => (
                  <SelectItem key={y.id} value={y.id.toString()}>
                    {y?.formatted}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
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
