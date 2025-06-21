import { DividerText } from "@/components/atoms/DividerText";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { MailCard } from "@/components/mails/MailCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { useMailGetAll } from "@/lib/api/mail";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Link, Outlet, useMatch } from "@tanstack/react-router";

export function Mails() {
  const isCreate = useMatch({ from: "/mails/create", shouldThrow: false });
  const isEdit = useMatch({ from: "/mails/edit/$mail", shouldThrow: false });

  const { data: mails, isLoading } = useMailGetAll()

  useBreadcrumb({ title: "Mails", link: { to: "/mails" } })

  if (isLoading) {
    return <Indeterminate />
  }

  if (!mails) {
    return null;
  }

  if (isCreate || isEdit) {
    return <Outlet />
  }

  const now = Date.now()
  const upcomingMails = mails.filter(m => m.sendTime.getTime() > now).sort((a, b) => a.sendTime.getTime() - b.sendTime.getTime())
  const passedMails = mails.filter(m => m.sendTime.getTime() <= now).sort((a, b) => b.sendTime.getTime() - a.sendTime.getTime())

  return (
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>Mails</Title>
        <Button size="lg" variant="outline" asChild>
          <Link to="/mails/create">
            Create
          </Link>
        </Button>
      </PageHeader>
      <div className="grid gap-4 grid-cols-1">
        {upcomingMails.map(m => (
          <MailCard key={m.id} mail={m} />
        ))}
      </div>
      {passedMails.length > 0 && (
        <>
          <DividerText>
            Past Announcements
          </DividerText>
          <div className="grid gap-4 grid-cols-1">
            {passedMails.map(m => (
              <MailCard key={m.id} mail={m} />
            ))}
          </div>
        </>
      )}
      {mails.length === 0 && (
        <div className="flex flex-col items-center space-y-4 pt-48">
          <h3 className="text-lg font-semibold">No mails found</h3>
          <h5 className="text-md text-muted-foreground">Get started by adding some</h5>
        </div>
      )}
    </div>
  )
}
