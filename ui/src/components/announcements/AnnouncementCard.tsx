import { useEventByYear } from "@/lib/api/event";
import { useYear } from "@/lib/hooks/useYear";
import { Announcement } from "@/lib/types/announcement";
import { formatDate } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { Badge } from "../ui/badge";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { Separator } from "../ui/separator";
import { Fragment } from "react/jsx-runtime";

interface Props {
  announcement: Announcement;
}

function AnnouncementBadge({ announcement }: { announcement: Announcement }) {
  if (announcement.send) {
    return <Badge variant="outline" className="text-green-500 border-green-500">Send</Badge>
  }

  if (announcement.error) {
    return <Badge variant="outline" className="text-red-500 border-red-500">Error</Badge>
  }

  return null
}

export function AnnouncementCard({ announcement }: Props) {
  const navigate = useNavigate()

  const { year } = useYear()
  const { data: allEvents, isLoading: isLoadingEvents } = useEventByYear(year)

  const events = allEvents?.filter(e => announcement.eventIds.includes(e.id)) ?? []

  if (isLoadingEvents) {
    return
  }

  const handleClick = () => {
    if (announcement.send || announcement.error) {
      return
    }

    navigate({ to: "/announcements/edit/$announcementId", params: { announcementId: announcement.id.toString() } })
  }

  return (
    <Card onClick={handleClick} className={!announcement.send && !announcement.error ? "transition-transform duration-300 cursor-pointer hover:scale-102" : ""}>
      <CardHeader>
        <div className="flex justify-between">
          <CardTitle>{formatDate(announcement.sendTime)}</CardTitle>
          <AnnouncementBadge announcement={announcement} />
        </div>
        <CardDescription>
          <div className="xs:flex md:grid md:grid-cols-[auto_1fr] md:space-x-2">
            {events.map(e => (
              <Fragment key={e.id}>
                <span>{e.name}</span>
                <span>{` | ${formatDate(e.startTime)}`}</span>
                <br className="md:hidden" />
              </Fragment>
            ))}
          </div>
        </CardDescription>
        <Separator />
      </CardHeader>
      <CardContent>
        <MarkdownViewer value={announcement.content} />
      </CardContent>
      {announcement.error && (
        <CardFooter>
          <span className="text-sm text-red-500">{announcement.error}</span>
        </CardFooter>
      )}
    </Card>
  )
}
