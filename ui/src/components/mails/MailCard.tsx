import { Mail } from "@/lib/types/mail";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { useNavigate } from "@tanstack/react-router";
import { formatDate } from "@/lib/utils/utils";
import { MarkdownViewer } from "../organisms/markdown/MarkdownViewer";

interface Props {
  mail: Mail
}

export function MailCard({ mail }: Props) {
  const navigate = useNavigate()

  const handleClick = () => {
    if (mail.send) {
      return
    }

    navigate({ to: "/mails/edit/$mail", params: { mail: mail.id.toString() } })
  }

  return (
    <Card onClick={handleClick} className={`transition-transform duration-300 hover:scale-101 ${!mail.send && "hover:cursor-pointer"}`}>
      <CardHeader>
        <div className="flex justify-between">
          <div>
            <CardTitle>{mail.title}</CardTitle>
            <CardDescription>
              <span>{`Send time: ${formatDate(mail.sendTime)}`}</span>
            </CardDescription>
          </div>
        </div>
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
