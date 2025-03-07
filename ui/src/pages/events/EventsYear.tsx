import { DividerText } from "@/components/atoms/DividerText";
import { EventCard } from "@/components/events/EventCard";
import { LoadingCards } from "@/components/organisms/LoadingCards";
import { Button } from "@/components/ui/button";
import { useEventByYear } from "@/lib/api/event";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Link, Outlet, useMatch, useParams } from "@tanstack/react-router";
import { isAfter, isBefore } from "date-fns";

export function EventsYear() {
  const isAssign = useMatch({ from: "/events/$year/assign", shouldThrow: false });
  const { year: yearString } = useParams({ from: "/events/$year" });

  const { data: years } = useYearGetAll();
  // Event component makes sure it exists
  const year = years!.find(({ formatted }) => formatted === yearString)!;
  const { data: events } = useEventByYear(year);

  useBreadcrumb({ title: yearString, link: { to: "/events/$year", params: { year: yearString } } });

  if (!events) {
    return <LoadingCards rows={4} cols={3} />;
  }

  if (isAssign) {
    return <Outlet />;
  }

  const now = Date.now();
  const futureEvents = events.filter(event => isAfter(event.endTime ?? event.startTime, now));
  const pastEvents = events.filter(event => isBefore(event.endTime ?? event.startTime, now)).sort((a, b) => b.startTime.getTime() - a.startTime.getTime());

  return (
    <div>
      <div className="flex pb-8 justify-between gap-6">
        Multiple Selector
        <div className="flex gap-6 items-center">
          <Button size="lg" variant="outline" asChild>
            <Link to="/events/$year/assign" params={{ year: yearString }}>
              Assign
            </Link>
          </Button>
          <Button>Sync TODO</Button>
        </div>
      </div>
      <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
        {futureEvents.map(event => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
      {pastEvents.length > 0 && (
        <>
          <DividerText text="Past Events" />
          <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
            {pastEvents.map(event => (
              <EventCard key={event.id} event={event} />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
