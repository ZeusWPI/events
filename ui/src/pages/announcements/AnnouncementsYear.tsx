import { AnnouncementCard } from "@/components/announcements/AnnouncementCard";
import { DividerText } from "@/components/atoms/DividerText";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useEventByYear } from "@/lib/api/event";
import { useYearGetAll } from "@/lib/api/year";
import { useIsMobile } from "@/lib/hooks/use-mobile";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Link, Outlet, useMatch, useParams } from "@tanstack/react-router";
import { ChevronRightIcon } from "lucide-react";

export function AnnouncementsYear() {
  const isCreate = useMatch({ from: "/announcements/$year/$event", shouldThrow: false });
  const { year: yearString } = useParams({ from: "/announcements/$year" });

  const { data: years } = useYearGetAll();
  // Announcements component makes sure it exists
  const year = years!.find(({ formatted }) => formatted === yearString)!;
  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year);

  useBreadcrumb({ title: yearString, link: { to: "/announcements/$year", params: { year: yearString } } });
  const isMobile = useIsMobile();

  if (isLoadingEvents) {
    return <Indeterminate />;
  }

  if (isCreate) {
    return <Outlet />
  }

  const now = Date.now()
  const announcements = events?.filter(e => e.announcement !== undefined) ?? []
  const noAnnouncements = events?.filter(e => e.announcement === undefined) ?? []
  const upcomingAnnouncements = announcements.filter(e => e.announcement!.sendTime.getTime() > now).sort((a, b) => a.announcement!.sendTime.getTime() - b.announcement!.sendTime.getTime())
  const passedAnnouncements = announcements.filter(e => e.announcement!.sendTime.getTime() <= now).sort((a, b) => b.announcement!.sendTime.getTime() - a.announcement!.sendTime.getTime())

  return (
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>{`${!isMobile ? "Announcements " : ""} ${yearString}`}</Title>
        <Dialog>
          <DialogTrigger asChild>
            <Button size="lg" variant="outline">
              Create
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                Select event to create an announcement for
              </DialogTitle>
              <DialogDescription asChild>
                <ScrollArea className="h-[600px] px-2">
                  <div className="divide-y">
                    {noAnnouncements.map(e => (
                      <div key={e.id} className="flex items-center justify-between py-2">
                        <span>{e.name}</span>
                        <Button size="icon" variant="ghost" asChild>
                          <Link to="/announcements/$year/$event" params={{ year: yearString, event: e.id.toString() }}>
                            <ChevronRightIcon />
                          </Link>
                        </Button>
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              </DialogDescription>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      </PageHeader>
      <div className="grid gap-4 grid-cols-1">
        {upcomingAnnouncements.map(e => (
          <AnnouncementCard key={e.id} event={e} />
        ))}
      </div>
      {passedAnnouncements.length > 0 && (
        <>
          <DividerText>
            Past Announcements
          </DividerText>
          <div className="grid gap-4 grid-cols-1">
            {passedAnnouncements.map(e => (
              <AnnouncementCard key={e.id} event={e} />
            ))}
          </div>
        </>
      )}
      {announcements.length === 0 && (
        <div className="flex flex-col items-center space-y-4 pt-48">
          <h3 className="text-lg font-semibold">No announcements found</h3>
          <h5 className="text-md text-muted-foreground">Get started by adding some</h5>
        </div>
      )}
    </div>
  );
}
