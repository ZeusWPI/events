import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { EventAssignCard } from "@/components/events/EventAssignCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useEventByYear, useEventSaveOrganizers } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { Link, useBlocker } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { ArrowLeft, LoaderCircle } from "lucide-react";
import { Fragment, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { DividerText } from "@/components/atoms/DividerText";
import { useIsMobile } from "@/lib/hooks/use-mobile";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { useYear, useYearLock } from "@/lib/hooks/useYear";
import { weightSubcategory } from "@/lib/types/general";
import { isAfter, isBefore } from "date-fns";

export function EventsAssign() {
  const { year } = useYear()

  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year);
  const { data: organizers, isLoading: isLoadingOrganizers } = useOrganizerByYear(year);

  const updateOrganizers = useEventSaveOrganizers();

  const [updatedEvents, setUpdatedEvents] = useState(events ?? []);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (!events) return
    setUpdatedEvents(events)
  }, [events])

  useBreadcrumb({ title: "Assign", weight: weightSubcategory, link: { to: "/events/assign" } });
  useYearLock()
  const isMobile = useIsMobile();

  const now = Date.now()

  const futureEvents = updatedEvents?.filter(event => isAfter(event.endTime ?? event.startTime, now)) ?? []
  const pastEvents = updatedEvents?.filter(event => isBefore(event.endTime ?? event.startTime, now)) ?? []

  const organizersEventCount = organizers?.map(organizer => ({
    ...organizer,
    futureEvents: updatedEvents.filter(event => isAfter(event.endTime ?? event.startTime, now) && event.organizers.some(({ id }) => id === organizer.id)).length,
    oldEvents: updatedEvents.filter(event => isBefore(event.endTime ?? event.startTime, now) && event.organizers.some(({ id }) => id === organizer.id)).length,
  })).sort((a, b) => b.futureEvents - a.futureEvents || a.name.localeCompare(b.name)) ?? []

  const handleSave = () => {
    setIsSaving(true);
    updateOrganizers.mutate(updatedEvents, {
      onSuccess: () => {
        toast.success("Success");
      },
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setIsSaving(false),
    });
  };

  const handleAssign = (eventId: number, selectedOrganizers: number[]) => {
    setUpdatedEvents(prev =>
      prev.map(event =>
        event.id === eventId
          ? { ...event, organizers: organizers?.filter(({ id }) => selectedOrganizers.includes(id)) ?? [] }
          : event,
      ),
    );
  };

  const changed = useMemo(() => JSON.stringify(events) !== JSON.stringify(updatedEvents), [events, updatedEvents])
  useBlocker({
    shouldBlockFn: () => {
      if (!changed) return false

      const shouldLeave = confirm("Are you sure you want to leave. You have unsaved changed.")
      return !shouldLeave
    }
  })

  if (isLoadingEvents || isLoadingOrganizers) {
    return <Indeterminate />;
  }

  return (
    <div className="grid xl:grid-cols-4 gap-8">
      <PageHeader className="col-span-full">
        <Title>{`Assign${!isMobile ? ` to Events ${year.formatted}` : ""}`}</Title>
        <div className="flex items-center gap-6">
          <Button size="lg" variant="outline" asChild>
            <Link to="/events">
              {changed
                ? "Discard"
                : (
                  <>
                    <ArrowLeft />
                    <span>Go back</span>
                  </>
                )}
            </Link>
          </Button>
          <Button disabled={!changed || isSaving} onClick={handleSave} className="w-16">
            {isSaving ? <LoaderCircle className="animate-spin" /> : "Save"}
          </Button>
        </div>
      </PageHeader>
      <div className="sticky top-6">
        <Card className={`w-full xl:max-w-80 sticky top-6 ${isMobile ? "max-h-48 overflow-y-scroll" : ""}`}>
          <CardContent className="flex flex-col gap-1 divide-y">
            {organizersEventCount.map(organizer => (
              <motion.div key={organizer.id} layout className="flex justify-between items-center w-full py-2">
                <span>{organizer.name}</span>
                <div>
                  <span>{organizer.futureEvents}</span>
                  <span className="text-muted-foreground ml-1">{`(${organizer.futureEvents + organizer.oldEvents})`}</span>
                </div>
              </motion.div>
            ))}
          </CardContent>
        </Card>
      </div>
      {futureEvents.length > 0 && (
        <div className="xl:col-span-3 flex flex-col gap-5">
          {futureEvents.map((event, i) => (
            <Fragment key={event.id}>
              <EventAssignCard event={event} organizers={organizers ?? []} onAssign={handleAssign} />
              {i !== (events?.length ?? 0) - 1 && <Separator />}
            </Fragment>
          ))}
        </div>
      )}
      {pastEvents.length > 0 && (
        <div className="xl:col-span-3 xl:col-start-2 flex flex-col gap-5">
          <DividerText>
            Past Events
          </DividerText>
          {pastEvents.map((event, i) => (
            <Fragment key={event.id}>
              <EventAssignCard event={event} organizers={organizers ?? []} onAssign={handleAssign} />
              {i !== (events?.length ?? 0) - 1 && <Separator />}
            </Fragment>
          ))}
        </div>
      )}
    </div>
  );
}
