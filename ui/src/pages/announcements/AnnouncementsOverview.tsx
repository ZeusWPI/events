import { AnnouncementCard } from "@/components/announcements/AnnouncementCard";
import { DividerText } from "@/components/atoms/DividerText";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { NoItems } from "@/components/atoms/NoItems";
import { Title } from "@/components/atoms/Title";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { useAnnouncementByYear } from "@/lib/api/announcement";
import { useIsMobile } from "@/lib/hooks/use-mobile";
import { useYear } from "@/lib/hooks/useYear";
import { Link } from "@tanstack/react-router";

export function AnnouncementsOverview() {
  const { year } = useYear()

  const { data: announcements, isLoading: isLoadingAnnouncements } = useAnnouncementByYear(year);

  const isMobile = useIsMobile();

  if (isLoadingAnnouncements) {
    return <Indeterminate />;
  }

  const now = Date.now()
  const upcomingAnnouncements = announcements?.filter(a => a.sendTime.getTime() > now).sort((a, b) => a.sendTime.getTime() - b.sendTime.getTime()) ?? []
  const passedAnnouncements = announcements?.filter(a => a.sendTime.getTime() <= now).sort((a, b) => b.sendTime.getTime() - a.sendTime.getTime()) ?? []

  return (
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>{`${!isMobile ? "Announcements " : ""} ${year.formatted}`}</Title>
        <Button size="lg" variant="outline" asChild>
          <Link to="/announcements/create">
            Create
          </Link>
        </Button>
      </PageHeader>
      <div className="grid gap-4 grid-cols-1">
        {upcomingAnnouncements.map(a => (
          <AnnouncementCard key={a.id} announcement={a} />
        ))}
      </div>
      {passedAnnouncements.length > 0 && (
        <>
          <DividerText>
            Past Announcements
          </DividerText>
          <div className="grid gap-4 grid-cols-1">
            {passedAnnouncements.map(a => (
              <AnnouncementCard key={a.id} announcement={a} />
            ))}
          </div>
        </>
      )}
      {announcements?.length === 0 && <NoItems title="No announcements found" description="Get started by clicking the create button" />}
    </div>
  );
}

// TODO: timedate selector
// Delete announcements
// Keep track of who made a mail / announcement
// Event multiple announcements and years
