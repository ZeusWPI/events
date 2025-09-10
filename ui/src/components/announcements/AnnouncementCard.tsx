import { useAnnouncementDelete, useAnnouncementResend } from "@/lib/api/announcement";
import { useEventByYear } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useYear } from "@/lib/hooks/useYear";
import { Announcement } from "@/lib/types/announcement";
import { formatDate } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { RotateCcwIcon, Trash2Icon } from "lucide-react";
import { useState } from "react";
import { Fragment } from "react/jsx-runtime";
import { toast } from "sonner";
import { Copy } from "../atoms/Copy";
import { IconButton } from "../atoms/IconButton";
import { OrganizerIcon } from "../atoms/OrganizerIcon";
import { TooltipText } from "../atoms/TooltipText";
import { DeleteConfirm } from "../molecules/DeleteConfirm";
import { ResendConfirm } from "../molecules/ResendConfirm";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { Badge } from "../ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Separator } from "../ui/separator";

interface Props {
  announcement: Announcement;
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

  const [openResend, setOpenResend] = useState(false)
  const announcementResend = useAnnouncementResend()

  if (isLoadingEvents || isLoadingOrganizers) {
    return
  }

  const handleClick = () => {
    if (announcement.send || announcement.error) {
      return
    }

    navigate({ to: "/announcements/edit/$announcementId", params: { announcementId: announcement.id.toString() } })
  }

  const handleResend = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    setOpenResend(true)
  }

  const handleResendConfirm = () => {
    announcementResend.mutate(announcement, {
      onSuccess: () => toast.success("Announcement resend", { description: "Scheduled for resending in one minute" }),
      onError: (err) => toast.error("Failed", { description: err.message }),
      onSettled: () => setOpenResend(false),
    })
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
              {organizer && <OrganizerIcon user={organizer} tooltip />}
              <CardTitle>{formatDate(announcement.sendTime)}</CardTitle>
              <AnnouncementBadge announcement={announcement} />
            </div>
            <ActionBar announcement={announcement} onResend={handleResend} onDelete={handleDelete} />
          </div>
          {announcement.error && <span className="text-sm text-red-500">{announcement.error}</span>}
          <CardDescription>
            <div className="xs:flex md:grid md:grid-cols-[auto_1fr] md:space-x-2">
              {events.map(e => (
                <Fragment key={e.id}>
                  <span>{`${e.name} | ${formatDate(e.startTime)}`}</span>
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
      </Card>
      <DeleteConfirm
        open={openDelete}
        onOpenChange={setOpenDelete}
        onDelete={handleDeleteConfirm}
      />
      <ResendConfirm
        open={openResend}
        onOpenChange={setOpenResend}
        onResend={handleResendConfirm}
      />
    </>
  )
}

interface ActionBarProps {
  announcement: Announcement;
  onResend: React.MouseEventHandler<HTMLButtonElement>;
  onDelete: React.MouseEventHandler<HTMLButtonElement>;
}

function ActionBar({ announcement, onResend, onDelete }: ActionBarProps) {
  return (
    <div className="flex items-center space-x-2">
      <Copy text={announcement.content} tooltip="Copy raw markdown" />
      {announcement.error && (
        <TooltipText text="Resend announcement">
          <IconButton onClick={onResend}>
            <RotateCcwIcon />
          </IconButton>
        </TooltipText>
      )}
      {!announcement.send && !announcement.error && (
        <TooltipText text="Delete">
          <IconButton onClick={onDelete}>
            <Trash2Icon className="text-red-500" />
          </IconButton>
        </TooltipText>
      )}
    </div>
  )
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


