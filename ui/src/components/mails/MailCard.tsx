import { useEventByYear } from "@/lib/api/event";
import { useMailDelete, useMailResend } from "@/lib/api/mail";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useYear } from "@/lib/hooks/useYear";
import { Mail } from "@/lib/types/mail";
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
import { Confirm } from "../molecules/Confirm";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { Badge } from "../ui/badge";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { Separator } from "../ui/separator";

interface Props {
  mail: Mail
}

export function MailCard({ mail }: Props) {
  const navigate = useNavigate()

  const { year } = useYear()
  const { data: allEvents, isLoading: isLoadingEvents } = useEventByYear(year)
  const { data: organizers, isLoading: isLoadingOrganizers } = useOrganizerByYear(year)

  const events = allEvents?.filter(e => mail.eventIds.includes(e.id)) ?? []
  const organizer = organizers?.find(o => o.id === mail.author_id)

  const [openDelete, setOpenDelete] = useState(false)
  const mailDelete = useMailDelete()

  const [openResend, setOpenResend] = useState(false)
  const mailResend = useMailResend()

  if (isLoadingEvents || isLoadingOrganizers) {
    return
  }

  const handleClick = () => {
    if (mail.send || mail.error) {
      return
    }

    navigate({ to: "/mails/edit/$mailId", params: { mailId: mail.id.toString() } })
  }

  const handleResend = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    setOpenResend(true)
  }

  const handleResendConfirm = () => {
    mailResend.mutate(mail, {
      onSuccess: () => toast.success("Mail resend", { description: "Scheduled for resending in one minute" }),
      onError: (err) => toast.error("Failed", { description: err.message }),
      onSettled: () => setOpenResend(false),
    })
  }

  const handleDelete = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    setOpenDelete(true)
  }

  const handleDeleteConfirm = () => {
    mailDelete.mutate(mail, {
      onSuccess: () => toast.success("Mail deleted"),
      onError: (err) => toast.error("Failed", { description: err.message }),
      onSettled: () => setOpenDelete(false),
    })
  }

  return (
    <>
      <Card onClick={handleClick} className={!mail.send && !mail.error ? "transition-transform duration-300 cursor-pointer hover:scale-102" : ""}>
        <CardHeader>
          <div className="flex justify-between">
            <CardTitle className="flex items-center space-x-2">
              {organizer && <OrganizerIcon user={organizer} tooltip />}
              <span>{mail.title}</span>
              <span className="font-normal text-muted-foreground text-sm">{` | ${formatDate(mail.sendTime)}`}</span>
              <MailBadge mail={mail} />
            </CardTitle>
            <ActionBar mail={mail} onResend={handleResend} onDelete={handleDelete} />
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
          <MarkdownViewer value={mail.content} />
        </CardContent>
        <CardFooter className="flex flex-col space-y-4 items-start">
          {mail.error && <span className="text-sm text-red-500">{mail.error}</span>}
        </CardFooter>
      </Card>
      <DeleteConfirm
        open={openDelete}
        onOpenChange={setOpenDelete}
        onDelete={handleDeleteConfirm}
      />
      <Confirm
        title="Resend confirmation"
        description="Are you sure you want to try to resend the mail"
        confirmText="Resend"
        onConfirm={handleResendConfirm}
        open={openResend}
        onOpenChange={setOpenResend}
      />
    </>
  )
}


interface ActionBarProps {
  mail: Mail;
  onResend: React.MouseEventHandler<HTMLButtonElement>;
  onDelete: React.MouseEventHandler<HTMLButtonElement>;
}

function ActionBar({ mail, onResend, onDelete }: ActionBarProps) {
  return (
    <div className="flex items-center space-x-2">
      <Copy text={mail.content} tooltip="Copy raw markdown" />
      {mail.error && (
        <TooltipText text="Resend mail">
          <IconButton onClick={onResend}>
            <RotateCcwIcon />
          </IconButton>
        </TooltipText>
      )}
      {!mail.send && !mail.error && (
        <TooltipText text="Delete">
          <IconButton onClick={onDelete}>
            <Trash2Icon className="text-red-500" />
          </IconButton>
        </TooltipText>
      )}
    </div>
  )
}

function MailBadge({ mail }: { mail: Mail }) {
  if (mail.send) {
    return <Badge variant="outline" className="text-green-500 border-green-500">Sent</Badge>
  }

  if (mail.error) {
    return <Badge variant="outline" className="text-red-500 border-red-500">Error</Badge>
  }

  return null
}
