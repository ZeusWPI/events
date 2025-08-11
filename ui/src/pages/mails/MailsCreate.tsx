import { MailForm } from "@/components/mails/MailForm";
import { useMailCreate } from "@/lib/api/mail";
import { useYearLock } from "@/lib/hooks/useYear";
import { MailSchema } from "@/lib/types/mail";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";

export function MailsCreate() {
  const navigate = useNavigate()
  const create = useMailCreate()

  useYearLock()

  const handleSubmit = (mail: MailSchema) => {
    create.mutate(mail, {
      onSuccess: () => {
        toast.success("Success")
        navigate({ to: "/mails" })
      },
      onError: error => toast.error("Failed", { description: error.message }),
    })

  }

  return (
    <MailForm onSubmit={handleSubmit} />
  )
}
