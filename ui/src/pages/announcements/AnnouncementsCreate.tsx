import { Indeterminate } from "@/components/atoms/Indeterminate";
import { useEventByYear } from "@/lib/api/event";
import { useYearGetAll } from "@/lib/api/year";
import { Link, useNavigate, useParams } from "@tanstack/react-router";
import Error404 from "../404";
import { PageHeader } from "@/components/molecules/PageHeader";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Title } from "@/components/atoms/Title";
import { Datalist, DatalistItem, DatalistItemTitle, DatalistItemContent } from "@/components/molecules/Datalist";
import { formatDate } from "@/lib/utils/utils";
import { useState } from "react";
import { HeadlessCard } from "@/components/molecules/HeadlessCard";
import { CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { DateTimePicker } from "@/components/organisms/DateTimePicker";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { useAnnouncementCreate, useAnnouncementUpdate } from "@/lib/api/announcement";
import { Announcement } from "@/lib/types/announcement";
import { MarkdownCombo } from "@/components/organisms/markdown/MarkdownCombo";

export function AnnouncementsCreate() {
  const { year: yearString, event: eventId } = useParams({ from: "/announcements/$year/$event" })

  const { data: years } = useYearGetAll();
  // Announcements component makes sure it exists
  const year = years!.find(({ formatted }) => formatted === yearString)!;
  const { data: events, isLoading } = useEventByYear(year);

  const event = events?.find(e => e.id === Number(eventId))
  const [date, setDate] = useState(event?.announcement?.sendTime)
  const [content, setContent] = useState(event?.announcement?.content)

  const announcementCreate = useAnnouncementCreate()
  const announcementUpdate = useAnnouncementUpdate()

  const navigate = useNavigate()

  useBreadcrumb({ title: "Create", link: { to: "/announcements/$year/$event", params: { year: yearString, event: eventId } } });

  if (isLoading) {
    return <Indeterminate />
  }

  if (!event) {
    return <Error404 />
  }

  const handleSubmit = () => {
    const now = Date.now()

    if (!date || date.getTime() <= now) {
      toast.error("Invalid date", { description: "Date has to be in the future" })
      return
    }
    if (event && event.startTime.getTime() < date.getTime()) {
      toast.error("Invalid date", { description: "Date has to be before the event" })
      return
    }

    if (!content || !content.length) {
      toast.error("Invalid announcement", { description: "Announcement text can't be empty" })
      return
    }

    const announcement: Announcement = {
      id: event.announcement?.id ?? 0,
      eventId: event.id,
      content,
      sendTime: date,
      send: event.announcement?.send ?? false,
    }

    let action
    if (announcement.id) {
      action = announcementUpdate
    } else {
      action = announcementCreate
    }

    action.mutate(announcement, {
      onSuccess: () => {
        toast.success("Success")
        navigate({ to: "/announcements/$year", params: { year: yearString } })
      },
      onError: error => toast.error("Failed", { description: error.message }),
    })
  }

  return (
    <div className="space-y-8">
      <PageHeader>
        <Title>{`${event.announcement?.content ? "Edit" : "Create"} Announcement`}</Title>
        <div className="flex justify-center gap-2">
          <Button variant="outline" asChild>
            <Link to="/announcements/$year" params={{ year: yearString }}>
              Cancel
            </Link>
          </Button>
          <Button onClick={handleSubmit} disabled={date?.getTime() === event.announcement?.sendTime.getTime() && content === event.announcement?.content}>
            Submit
          </Button>
        </div>
      </PageHeader>
      <HeadlessCard>
        <CardHeader className="px-4 sm:px-0 pt-0">
          <CardTitle>{event.name}</CardTitle>
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
      <HeadlessCard>
        <CardHeader className="px-4 sm:px-0 pt-0">
          <CardTitle>Announcement</CardTitle>
        </CardHeader>
        <CardContent className="px-0 space-y-4">
          <div className="flex items-center gap-4">
            <span>Send time</span>
            <DateTimePicker value={date} onChange={setDate} weekStartsOn={1} className="w-[280px]" />
          </div>
          <MarkdownCombo value={content} onChange={setContent} textareaProps={{ placeholder: "Write announcement here..." }} />
        </CardContent>
      </HeadlessCard>
    </div>
  )
}
