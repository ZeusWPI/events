import { useAnnouncementDelete } from "@/lib/api/announcement";
import { useEventByYear } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useYear } from "@/lib/hooks/useYear";
import { Announcement } from "@/lib/types/announcement";
import { formatDate } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { Trash2Icon } from "lucide-react";
import { useState } from "react";
import { Fragment } from "react/jsx-runtime";
import { toast } from "sonner";
import { OrganizerIcon } from "../atoms/OrganizerIcon";
import { DeleteConfirm } from "../molecules/DeleteConfirm";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { Separator } from "../ui/separator";

interface Props {
  announcement: Announcement;
}

function AnnouncementBadge({ announcement, onDelete }: { announcement: Announcement, onDelete: React.MouseEventHandler<HTMLButtonElement> }) {
  if (announcement.send) {
    return <Badge variant="outline" className="text-green-500 border-green-500">Send</Badge>
  }

  if (announcement.error) {
    return <Badge variant="outline" className="text-red-500 border-red-500">Error</Badge>
  }

  return (
    <Button onClick={onDelete} size="icon" variant="secondary" className="size-6">
      <Trash2Icon className="text-red-500" />
    </Button>
  )
}

export function AnnouncementCard({ announcement }: Props) {
  const navigate = useNavigate()

  const { year } = useYear()
  const { data: allEvents, isLoading: isLoadingEvents } = useEventByYear(year)
  const { data: organizers, isLoading: isLoadingOrganizers } = useOrganizerByYear(year)

  const events = allEvents?.filter(e => announcement.eventIds.includes(e.id)) ?? []
  const organizer = organizers?.find(o => o.id === announcement.author_id)

  const [openDelete, setOpenDelete] = useState(false)
  const announcementDelete = useAnnouncementDelete()

  if (isLoadingEvents || isLoadingOrganizers) {
    return
  }

  const handleClick = () => {
    if (announcement.send || announcement.error) {
      return
    }

    navigate({ to: "/announcements/edit/$announcementId", params: { announcementId: announcement.id.toString() } })
  }

  const handleDelete = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    setOpenDelete(true)
  }

  const handleDeleteConfirm = () => {
    announcementDelete.mutate(announcement, {
      onSuccess: () => toast.success("Announcement deleted"),
      onError: (err) => toast.error("Failed", { description: err.message }),
      onSettled: () => setOpenDelete(false),
    })
  }

  return (
    <>
      <Card onClick={handleClick} className={!announcement.send && !announcement.error ? "transition-transform duration-300 cursor-pointer hover:scale-102" : ""}>
        <CardHeader>
          <div className="flex justify-between">
            <div className="flex items-center space-x-2">
              {organizer && <OrganizerIcon user={organizer} />}
              <CardTitle>{formatDate(announcement.sendTime)}</CardTitle>
            </div>
            <AnnouncementBadge announcement={announcement} onDelete={handleDelete} />
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
        <CardFooter className="flex flex-col space-y-4 items-start">
          {announcement.error && <span className="text-sm text-red-500">{announcement.error}</span>}
        </CardFooter>
      </Card>
      <DeleteConfirm
        open={openDelete}
        onOpenChange={setOpenDelete}
        onDelete={handleDeleteConfirm}
      />
    </>
  )
}
