import { AnnouncementList } from "@/components/announcements/AnnouncementList";
import { ButtonGroup } from "@/components/atoms/ButtonGroup";
import { IconButton } from "@/components/atoms/IconButton";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { CheckTable } from "@/components/check/CheckTable";
import { EventPoster } from "@/components/events/EventPoster";
import { MailList } from "@/components/mails/MailList";
import { Datalist, DatalistItem, DatalistItemContent, DatalistItemTitle } from "@/components/molecules/Datalist";
import { HeadlessCard } from "@/components/molecules/HeadlessCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useEventByYear } from "@/lib/api/event";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { useYear, useYearLock } from "@/lib/hooks/useYear";
import { weightSubcategory } from "@/lib/types/general";
import { formatDate } from "@/lib/utils/utils";
import { Link, useParams } from "@tanstack/react-router";
import { ChevronDownIcon, ChevronUpIcon, LinkIcon, MailIcon, MegaphoneIcon, UserRound } from "lucide-react";
import { useState } from "react";
import Error404 from "../404";

export function EventsDetail() {
  const { id: eventID } = useParams({ from: "/events/$id" });
  const { year } = useYear()

  const { data: events, isLoading } = useEventByYear(year);
  const event = events?.find(event => event.id.toString() === eventID);

  const big = event?.posters.find(p => !p.scc) ?? { id: 0, eventId: event?.id ?? 0, scc: false }
  const scc = event?.posters.find(p => p.scc) ?? { id: 0, eventId: event?.id ?? 0, scc: true }

  const [showAnnouncements, setShowAnnouncements] = useState(false)
  const [showMails, setShowMails] = useState(false)

  useBreadcrumb({ title: event?.name ?? "", weight: weightSubcategory, link: { to: "/events/$id", params: { id: eventID } } });
  useYearLock()

  if (isLoading) {
    return <Indeterminate />
  }

  if (!event) {
    return <Error404 />;
  }

  return (
    <div className="grid lg:grid-cols-6 gap-8">
      <PageHeader className="col-span-full">
        <Title>{event.name}</Title>
        <ButtonGroup>
          <IconButton>
            <Link to="/announcements/create">
              <MegaphoneIcon />
            </Link>
          </IconButton>
          <IconButton>
            <Link to="/mails/create">
              <MailIcon />
            </Link>
          </IconButton>
          <IconButton asChild>
            <a href={event.url} rel="noopener noreferrer" target="_blank">
              <LinkIcon />
            </a>
          </IconButton>
        </ButtonGroup>
      </PageHeader>
      <div className="lg:col-span-4">
        <HeadlessCard>
          <CardHeader className="px-4 sm:px-0 pt-0">
            <CardTitle>General</CardTitle>
          </CardHeader>
          <CardContent className="px-0">
            <Datalist>
              <DatalistItem>
                <DatalistItemTitle>Description</DatalistItemTitle>
                <DatalistItemContent>{event.description}</DatalistItemContent>
              </DatalistItem>
              <DatalistItem>
                <DatalistItemTitle>Time</DatalistItemTitle>
                <DatalistItemContent>
                  <div className="flex flex-col">
                    <span>{formatDate(event.startTime)}</span>
                    {event.endTime && (
                      <span>{formatDate(event.endTime)}</span>
                    )}
                  </div>
                </DatalistItemContent>
              </DatalistItem>
              <DatalistItem>
                <DatalistItemTitle>Location</DatalistItemTitle>
                <DatalistItemContent>{event.location}</DatalistItemContent>
              </DatalistItem>
            </Datalist>
          </CardContent>
        </HeadlessCard>
      </div>
      <div className="lg:col-span-2">
        <div className="flex h-full">
          <Separator orientation="vertical" className="h-full hidden lg:block" />
          <HeadlessCard className="h-full mx-auto">
            <CardHeader className="px-0 pt-0">
              <CardTitle className="flex gap-2 items-center">
                <UserRound />
                Organizers
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col space-y-1 px-0">
              {event.organizers.length
                ? event.organizers.map(organizer => (
                  <span key={organizer.id}>{organizer.name}</span>
                ))
                : (
                  <span>No one assigned</span>
                )}
            </CardContent>
          </HeadlessCard>
        </div>
      </div>
      <div className="col-span-full">
        <CheckTable checks={event.checks} eventId={Number(eventID)} />
      </div>
      <div className="col-span-full py-4">
        <div className="flex flex-wrap gap-16 justify-center lg:justify-start">
          <EventPoster
            title="Big Poster"
            description="Poster to distribute"
            poster={big}
            year={year}
          />
          <EventPoster
            title="SCC Poster"
            description="Cammiechat screen poster"
            poster={scc}
            year={year}
          />
        </div>
      </div>
      {event.announcements.length > 0 && (
        <div className="lg:col-span-3">
          <HeadlessCard>
            <CardHeader className="px-4 sm:px-0 pt-0">
              <CardTitle className="flex justify-between items-center">
                <span>{`${event.announcements.length} Announcement${event.announcements.length !== 1 ? 's' : ''}`}</span>
                <IconButton onClick={() => setShowAnnouncements(prev => !prev)} disabled={event.announcements.length === 0}>
                  {showAnnouncements
                    ? <ChevronUpIcon />
                    : <ChevronDownIcon />
                  }
                </IconButton>
              </CardTitle>
            </CardHeader>
            {showAnnouncements && (
              <CardContent className="px-0">
                <AnnouncementList announcements={event.announcements} />
              </CardContent>
            )}
          </HeadlessCard>
        </div>
      )}
      {event.mails.length > 0 && (
        <div className="lg:col-span-3">
          <HeadlessCard>
            <CardHeader className="px-4 sm:px-0 pt-0">
              <CardTitle className="flex justify-between items-center">
                <span>{`${event.mails.length} Mail${event.mails.length !== 1 ? 's' : ''}`}</span>
                <IconButton onClick={() => setShowMails(prev => !prev)} disabled={event.announcements.length === 0}>
                  {showMails
                    ? <ChevronUpIcon />
                    : <ChevronDownIcon />
                  }
                </IconButton>
              </CardTitle>
            </CardHeader>
            {showMails && (
              <CardContent className="px-0">
                <MailList mails={event.mails} />
              </CardContent>
            )}
          </HeadlessCard>
        </div>
      )}
    </div>
  );
}
