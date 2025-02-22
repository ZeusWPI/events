import type { ReactNode } from "react";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarHeader, SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { useBreadcrumbs } from "@/lib/hooks/useBreadcrumb";
import { Link } from "@tanstack/react-router";
import { Fragment } from "react";
import NavEvents from "./NavEvents";

function AppSidebar({ children }: { children: ReactNode }) {
  const { state: breadcrumbs } = useBreadcrumbs();

  return (
    <SidebarProvider>
      <Sidebar variant="inset" collapsible="icon">
        <SidebarHeader>
          Header
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <NavEvents />
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          Footer
        </SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <div className="container mx-auto px-2 pt-1 h-full">
          <header className="flex h-16 shrink-0 items-center gap-2">
            <SidebarTrigger />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                {breadcrumbs.map((crumb, index) => (
                  <Fragment key={crumb.title}>
                    <BreadcrumbItem>
                      {crumb.link
                        ? (
                            <BreadcrumbLink>
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
