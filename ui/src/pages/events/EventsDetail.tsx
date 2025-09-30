import { AnnouncementList } from "@/components/announcements/AnnouncementList";
import { ButtonGroup } from "@/components/atoms/ButtonGroup";
import { IconButton } from "@/components/atoms/IconButton";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { TooltipText } from "@/components/atoms/TooltipText";
import { CheckTable } from "@/components/check/CheckTable";
import { EventPoster } from "@/components/events/EventPoster";
import { MailList } from "@/components/mails/MailList";
import { Datalist, DatalistItem, DatalistItemContent, DatalistItemTitle } from "@/components/molecules/Datalist";
import { HeadlessCard } from "@/components/molecules/HeadlessCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useEvent } from "@/lib/api/event";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { useYear } from "@/lib/hooks/useYear";
import { weightSubcategory } from "@/lib/types/general";
import { formatDate } from "@/lib/utils/utils";
import { Link, useParams } from "@tanstack/react-router";
import { ChevronDownIcon, ChevronUpIcon, LinkIcon, MailIcon, MegaphoneIcon, UserRound } from "lucide-react";
import { useState } from "react";
import Error404 from "../404";
import { OrganizerIcon } from "@/components/atoms/OrganizerIcon";

export function EventsDetail() {
  const { id: eventID } = useParams({ from: "/events/$id" });
  const { year } = useYear()

  const { data: event, isLoading } = useEvent({ id: Number(eventID) });

  const big = event?.posters.find(p => !p.scc)
  const scc = event?.posters.find(p => p.scc)

  const [showAnnouncements, setShowAnnouncements] = useState(false)
  const [showMails, setShowMails] = useState(false)

  useBreadcrumb({ title: event?.name ?? "", weight: weightSubcategory, link: { to: "/events/$id", params: { id: eventID } } });

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
          <TooltipText text="Create announcement">
            <IconButton asChild>
              <Link to="/announcements/create">
                <MegaphoneIcon />
              </Link>
            </IconButton>
          </TooltipText>
          <TooltipText text="Create mail">
            <IconButton asChild>
              <Link to="/mails/create">
                <MailIcon />
              </Link>
            </IconButton>
          </TooltipText>
          <TooltipText text="Go to event website page">
            <IconButton asChild>
              <a href={event.url} rel="noopener noreferrer" target="_blank">
                <LinkIcon />
              </a>
            </IconButton>
          </TooltipText>
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
              <CardTitle className="flex gap-2 justify-center items-center">
                <UserRound />
                Organizers
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col space-y-1 px-0">
              {event.organizers.length
                ? (
                  <div className="table border-spacing-x-4 border-spacing-y-1">
                    {event.organizers.map(organizer => (
                      <div key={organizer.id} className="table-row">
                        <div className="table-cell">
                          <OrganizerIcon user={organizer} />
                        </div>
                        <div className="table-cell align-middle">{organizer.name}</div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <span>No one assigned</span>
                )}
            </CardContent>
          </HeadlessCard>
        </div>
      </div>
      <div className="col-span-full">
        <CheckTable checks={event.checks} event={event} />
      </div>
      <div className="lg:col-span-2 xl:col-span-2 w-full">
        <EventPoster
          title="Big Poster"
          description="Poster to distribute"
          poster={big}
          eventId={event.id}
          scc={false}
          year={year}
        />
      </div>
      <div className="lg:col-span-2 xl:col-span-2 w-full">
        <EventPoster
          title="SCC Poster"
          description="Cammiechat screen poster"
          poster={scc}
          eventId={event.id}
          scc={true}
          year={year}
        />
      </div>
      <div className="lg:col-span-full xl:col-span-2 flex flex-col">
        {event.announcements.length > 0 && (
          <div className="">
            <HeadlessCard>
              <CardHeader className="px-4 sm:px-0 pt-0">
                <CardTitle onClick={() => setShowAnnouncements(prev => !prev)} className="flex justify-between items-center cursor-pointer pb-1 border-b">
                  <span>{`${event.announcements.length} Announcement${event.announcements.length !== 1 ? 's' : ''}`}</span>
                  <Button size="icon" variant="ghost" disabled={event.announcements.length === 0}>
                    {showAnnouncements
                      ? <ChevronUpIcon />
                      : <ChevronDownIcon />
                    }
                  </Button>
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
          <div className="">
            <HeadlessCard>
              <CardHeader className="px-4 sm:px-0 pt-0">
                <CardTitle onClick={() => setShowMails(prev => !prev)} className="flex justify-between items-center cursor-pointer pb-1 border-b">
                  <span>{`${event.mails.length} Mail${event.mails.length !== 1 ? 's' : ''}`}</span>
                  <Button size="icon" variant="ghost" disabled={event.announcements.length === 0}>
                    {showMails
                      ? <ChevronUpIcon />
                      : <ChevronDownIcon />
                    }
                  </Button>
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
    </div>
  );
}
