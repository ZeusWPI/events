import { Mail } from "@/lib/types/mail";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { useNavigate } from "@tanstack/react-router";
import { formatDate } from "@/lib/utils/utils";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";
import { useEventByYear } from "@/lib/api/event";
import { useYear } from "@/lib/hooks/useYear";
import { Fragment } from "react/jsx-runtime";
import { Separator } from "../ui/separator";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { Trash2Icon } from "lucide-react";
import { useMailDelete } from "@/lib/api/mail";
import { useState } from "react";
import { toast } from "sonner";
import { DeleteConfirm } from "../molecules/DeleteConfirm";

interface Props {
  mail: Mail
}

function MailBadge({ mail, onDelete }: { mail: Mail, onDelete: React.MouseEventHandler<HTMLButtonElement> }) {
  if (mail.send) {
    return <Badge variant="outline" className="text-green-500 border-green-500">Send</Badge>
  }

  if (mail.error) {
    return <Badge variant="outline" className="text-red-500 border-red-500">Error</Badge>
  }

  return (
    <Button onClick={onDelete} size="icon" variant="secondary" className="size-6">
      <Trash2Icon className="text-red-500" />
    </Button>
  )
}

export function MailCard({ mail }: Props) {
  const navigate = useNavigate()

  const { year } = useYear()
  const { data: allEvents, isLoading: isLoadingEvents } = useEventByYear(year)

  const events = allEvents?.filter(e => mail.eventIds.includes(e.id)) ?? []

  const [openDelete, setOpenDelete] = useState(false)
  const mailDelete = useMailDelete()

  if (isLoadingEvents) {
    return
  }

  const handleClick = () => {
    if (mail.send || mail.error) {
      return
    }

    navigate({ to: "/mails/edit/$mailId", params: { mailId: mail.id.toString() } })
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
            <CardTitle>
              <span>{mail.title}</span>
              <span className="font-normal text-muted-foreground text-sm">{` | ${formatDate(mail.sendTime)}`}</span>
            </CardTitle>
            <MailBadge mail={mail} onDelete={handleDelete} />
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
        {mail.error && (
          <CardFooter>
            <span className="text-sm text-red-500">{mail.error}</span>
          </CardFooter>
        )}
      </Card>
      <DeleteConfirm
        open={openDelete}
        onOpenChange={setOpenDelete}
        onDelete={handleDeleteConfirm}
      />
    </>
  )
}
