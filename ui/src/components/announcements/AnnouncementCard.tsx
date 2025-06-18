import { Event } from "@/lib/types/event";
import { formatDate, formatDateDiff } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { Separator } from "../ui/separator";

interface Props {
  event: Event;
}

export function AnnouncementCard({ event }: Props) {
  const navigate = useNavigate()

  if (!event.announcement) {
    return null
  }

  const handleClick = () => {
    if (event.announcement!.send) {
      return
    }

    navigate({ to: "/announcements/$year/$event", params: { year: event.year.formatted, event: event.id.toString() } })
  }

  return (
    <Card onClick={handleClick} className={`transition-transform duration-300 hover:scale-101 ${!event.announcement.send && "hover:cursor-pointer"}`}>
      <CardHeader>
        <CardTitle>{event.name}</CardTitle>
        <CardDescription>
          <div className="flex flex-col">
            <span>{`Announcement: ${formatDate(event.announcement.sendTime)}`}</span>
            <span>{`Event: ${formatDate(event.startTime)} | ${formatDateDiff(event.announcement.sendTime, event.startTime)}`}</span>
          </div>
        </CardDescription>
        <Separator />
      </CardHeader>
      <CardContent>
        <MarkdownViewer value={event.announcement.content} />
      </CardContent>
    </Card>
  )
}
