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

interface Props {
  mail: Mail
}

function MailBadge({ mail }: { mail: Mail }) {
  if (mail.send) {
    return <Badge variant="outline" className="text-green-500 border-green-500">Send</Badge>
  }

  if (mail.error) {
    return <Badge variant="outline" className="text-red-500 border-red-500">Error</Badge>
  }

  return null
}

export function MailCard({ mail }: Props) {
  const navigate = useNavigate()

  const { year } = useYear()
  const { data: allEvents, isLoading: isLoadingEvents } = useEventByYear(year)

  const events = allEvents?.filter(e => mail.eventIds.includes(e.id)) ?? []

  if (isLoadingEvents) {
    return
  }

  const handleClick = () => {
    if (mail.send || mail.error) {
      return
    }

    navigate({ to: "/mails/edit/$mailId", params: { mailId: mail.id.toString() } })
  }

  return (
    <Card onClick={handleClick} className={!mail.send && !mail.error ? "transition-transform duration-300 cursor-pointer hover:scale-102" : ""}>
      <CardHeader>
        <div className="flex justify-between">
          <CardTitle>
            <span>{mail.title}</span>
            <span className="font-normal text-muted-foreground text-sm">{` | ${formatDate(mail.sendTime)}`}</span>
          </CardTitle>
          <MailBadge mail={mail} />
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
  )
}
