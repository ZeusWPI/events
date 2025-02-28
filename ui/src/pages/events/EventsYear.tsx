import { DividerText } from "@/components/atoms/DividerText";
import { EventCard } from "@/components/events/EventCard";
import { LoadingCards } from "@/components/molecules/LoadingCards";
import { Button } from "@/components/ui/button";
import { useEventByYear } from "@/lib/api/event";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { yearToString } from "@/lib/utils/converter";
import { useParams } from "@tanstack/react-router";
import { isAfter, isBefore } from "date-fns";

export function EventsYear() {
  const { yearId } = useParams({ from: "/events/$yearId" });

  const { data: years } = useYearGetAll();
  // Event component makes sure it exists
  const year = years!.find(({ id }) => id.toString() === yearId)!;
  const { data: events } = useEventByYear(year);

  useBreadcrumb({ title: yearToString(year), link: { to: "/events/$yearId", params: { yearId } } });

  if (!events) {
    return <LoadingCards rows={4} cols={3} />;
  }

  const now = Date.now();
  const futureEvents = events.filter(event => isAfter(event.endTime ?? event.startTime, now));
  const pastEvents = events.filter(event => isBefore(event.endTime ?? event.startTime, now)).sort((a, b) => b.startTime.getTime() - a.startTime.getTime());

  return (
    <div>
      <div className="flex pb-4 justify-between">
        <Button variant="outline">{yearToString(year)}</Button>
        <div className="space-x-6">
          <Button size="lg" variant="outline">Assign</Button>
          <Button size="lg">Sync</Button>
        </div>
      </div>
      <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
        {futureEvents.map(event => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
      <DividerText text="Past Events" />
      <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
        {pastEvents.map(event => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
    </div>
  );
}
