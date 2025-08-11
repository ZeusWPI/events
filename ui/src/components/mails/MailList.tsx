import { Mail } from "@/lib/types/mail";
import { MailCard } from "./MailCard";

interface Props {
  mails: Mail[];
}

export function MailList({ mails }: Props) {
  return (
    <div className="grid gap-4 grid-cols-1">
      {mails.map(m => <MailCard key={m.id} mail={m} />)}
    </div>
  )
}
