import { DividerText } from "@/components/atoms/DividerText";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { NoItems } from "@/components/atoms/NoItems";
import { Title } from "@/components/atoms/Title";
import { MailCard } from "@/components/mails/MailCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { useMailByYear } from "@/lib/api/mail";
import { useIsMobile } from "@/lib/hooks/use-mobile";
import { useYear } from "@/lib/hooks/useYear";
import { Link } from "@tanstack/react-router";

export function MailsOverview() {
  const { year } = useYear()

  const { data: mails, isLoading: isLoadingMails } = useMailByYear(year);

  const isMobile = useIsMobile();

  if (isLoadingMails) {
    return <Indeterminate />;
  }

  const now = Date.now()
  const upcomingMails = mails?.filter(a => a.sendTime.getTime() > now).sort((a, b) => a.sendTime.getTime() - b.sendTime.getTime()) ?? []
  const passedMails = mails?.filter(a => a.sendTime.getTime() <= now).sort((a, b) => b.sendTime.getTime() - a.sendTime.getTime()) ?? []

  return (
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>{`${!isMobile ? "Mails " : ""} ${year.formatted}`}</Title>
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
            Past Mails
          </DividerText>
          <div className="grid gap-4 grid-cols-1">
            {passedMails.map(m => (
              <MailCard key={m.id} mail={m} />
            ))}
          </div>
        </>
      )}
      {mails?.length === 0 && <NoItems title="No mails found" description="Get started by clicking the create button" />}
    </div>
  );
}
