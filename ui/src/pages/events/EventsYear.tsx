import { DividerText } from "@/components/atoms/DividerText";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { EventCard } from "@/components/events/EventCard";
import { MultiSelect } from "@/components/organisms/MultiSelect";
import { Button } from "@/components/ui/button";
import { useEventByYear } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Link, Outlet, useMatch, useParams } from "@tanstack/react-router";
import { isAfter, isBefore } from "date-fns";
import { useEffect, useState } from "react";

export function EventsYear() {
  const isDetail = useMatch({ from: "/events/$year/$id", shouldThrow: false });
  const isAssign = useMatch({ from: "/events/$year/assign", shouldThrow: false });
  const { year: yearString } = useParams({ from: "/events/$year" });

  const { data: years } = useYearGetAll();
  // Event component makes sure it exists
  const year = years!.find(({ formatted }) => formatted === yearString)!;
  const { data: events } = useEventByYear(year);

  const { data: organizers } = useOrganizerByYear(year);
  const [selectedOrganizers, setSelectedOrganizers] = useState<string[]>([]);
  useEffect(() => setSelectedOrganizers([]), [yearString]);

  useBreadcrumb({ title: yearString, link: { to: "/events/$year", params: { year: yearString } } });

  if (!events) {
    return <Indeterminate />;
  }

  if (isDetail || isAssign) {
    return <Outlet />;
  }

  const handleValueChange = (value: string[]) => setSelectedOrganizers(value);

  const now = Date.now();
  const filteredEvents = selectedOrganizers.length ? events.filter(event => event.organizers.find(({ id }) => selectedOrganizers.includes(id.toString()))) : events;
  const futureEvents = filteredEvents.filter(event => isAfter(event.endTime ?? event.startTime, now));
  const pastEvents = filteredEvents.filter(event => isBefore(event.endTime ?? event.startTime, now)).sort((a, b) => b.startTime.getTime() - a.startTime.getTime());

  return (
    <div>
      <div className="flex pb-8 justify-between gap-6">
        <div className="flex items-center gap-4">
          <MultiSelect
            options={organizers?.map(({ id, name }) => ({ value: id.toString(), label: name })) ?? []}
            onValueChange={handleValueChange}
            defaultValue={selectedOrganizers}
            placeholder="Filter by organizer"
            maxCount={2}
          />
        </div>
        <div className="flex gap-6 items-center">
          <Button size="lg" variant="outline" asChild>
            <Link to="/events/$year/assign" params={{ year: yearString }}>
              Assign
            </Link>
          </Button>
          {/* TODO: Implement */}
          <Button disabled>Sync</Button>
        </div>
      </div>
      <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
        {futureEvents.map(event => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
      {pastEvents.length > 0 && (
        <>
          <DividerText>
            Past Events
          </DividerText>
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
