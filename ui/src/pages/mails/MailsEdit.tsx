import { Indeterminate } from "@/components/atoms/Indeterminate"
import { MailForm } from "@/components/mails/MailForm"
import { useMailByYear, useMailUpdate } from "@/lib/api/mail"
import { useYear, useYearLock } from "@/lib/hooks/useYear"
import { MailSchema } from "@/lib/types/mail"
import { useParams, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"
import Error404 from "../404"

export function MailsEdit() {
  const { mailId } = useParams({ from: "/mails/edit/$mailId" })

  const { year } = useYear()
  const { data: mails, isLoading: isLoadingMails } = useMailByYear(year)
  const mail = mails?.find(a => a.id === Number(mailId))

  useYearLock()
  const navigate = useNavigate()
  const update = useMailUpdate()

  if (isLoadingMails) {
    return <Indeterminate />
  }

  if (!mail) {
    return <Error404 />
  }

  const handleSubmit = (mail: MailSchema) => {
    update.mutate(mail, {
      onSuccess: () => {
        toast.success("Success")
        navigate({ to: "/mails" })
      },
      onError: error => toast.error("Failed", { description: error.message }),
    })
  }

  return (
    <MailForm mail={mail} onSubmit={handleSubmit} />
  )
}
