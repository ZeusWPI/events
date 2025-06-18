import { Event } from "@/lib/types/event";
import { formatDate, formatDateDiff } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { Separator } from "../ui/separator";
import { Badge } from "../ui/badge";

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

  console.log(event.announcement)

  return (
    <Card onClick={handleClick} className={`transition-transform duration-300 hover:scale-101 ${!event.announcement.send && "hover:cursor-pointer"}`}>
      <CardHeader>
        <div className="flex justify-between">
          <div>
            <CardTitle>{event.name}</CardTitle>
            <CardDescription>
              <div className="flex flex-col">
                <span>{`Announcement: ${formatDate(event.announcement.sendTime)}`}</span>
                <span>{`Event: ${formatDate(event.startTime)} | ${formatDateDiff(event.announcement.sendTime, event.startTime)}`}</span>
              </div>
            </CardDescription>
          </div>
          {(event.announcement.send || event.announcement.error) && (
            <div>
              <Badge variant="outline" className={event.announcement.send ? "text-green-500 border-green-500" : "text-red-500 border-red-500"}>{event.announcement.send ? "Send" : "Error"}</Badge>
            </div>
          )}
        </div>
        <Separator />
      </CardHeader>
      <CardContent>
        <MarkdownViewer value={event.announcement.content} />
      </CardContent>
      {event.announcement.error && (
        <CardFooter>
          <span className="text-sm text-red-500">{event.announcement.error}</span>
        </CardFooter>
      )}
    </Card>
  )
}
