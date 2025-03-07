import { EventAssignCard } from "@/components/events/EventAssignCard";
import { LoadingCards } from "@/components/organisms/LoadingCards";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useEventByYear, useEventSaveOrganizers } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Link, useParams } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { ArrowLeft, LoaderCircle } from "lucide-react";

import { Fragment, useState } from "react";
import { toast } from "sonner";

export function EventsAssign() {
  const { year: yearString } = useParams({ from: "/events/$year" });

  const { data: years } = useYearGetAll();
  // Event component makes sure it exists
  const year = years!.find(({ formatted }) => formatted === yearString)!;
  const { data: events } = useEventByYear(year);
  const { data: organizers } = useOrganizerByYear(year);

  const updateOrganizers = useEventSaveOrganizers();

  const [updatedEvents, setUpdatedEvents] = useState(events ?? []);
  const [isDirty, setIsDirty] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  useBreadcrumb({ title: "Assign", link: { to: "/events/$year/assign", params: { year: yearString } } });

  const organizersEventCount = organizers?.map(organizer => ({
    ...organizer,
    events: updatedEvents.filter(event => event.organizers.some(({ id }) => id === organizer.id)).length,
  })).sort((a, b) => b.events - a.events || a.name.localeCompare(b.name)) ?? [];

  const handleDiscard = () => {
    // TODO: Ask for confirmation ?
    setIsDirty(false);
  };

  const handleSave = () => {
    setIsSaving(true);
    updateOrganizers.mutate(updatedEvents, {
      onSuccess: () => {
        setIsDirty(false);
        toast.success("🥳");
      },
      onError: error => toast.error("💀", { description: error.message }),
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
    setIsDirty(true);
  };

  if (!events) {
    return <LoadingCards rows={4} cols={3} />;
  }

  return (
    <div>
      <div className="flex pb-8 justify-end gap-6 items-center">
        <Button size="lg" variant="outline" onClick={handleDiscard} asChild>
          <Link to="/events/$year" params={{ year: yearString }}>
            {isDirty
              ? "Discard"
              : (
                  <>
                    <ArrowLeft />
                    <span>Go back</span>
                  </>
                )}
          </Link>
        </Button>
        <Button disabled={!isDirty || isSaving} onClick={handleSave} className="w-16">
          {isSaving ? <LoaderCircle className="animate-spin" /> : "Save"}
        </Button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-4">
        <div className="flex flex-col items-start border-muted-foreground">
          <Card className="w-full max-w-60 sticky top-6">
            <CardHeader>
              <CardTitle>User assignments</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-1">
              {organizersEventCount.map(organizer => (
                <motion.div key={organizer.id} layout className="flex justify-between items-center w-full max-w-60">
                  <span>{organizer.name}</span>
                  <span>{organizer.events}</span>
                </motion.div>
              ))}
            </CardContent>
          </Card>
        </div>
        <div className="md:col-span-3">
          <div className="flex flex-col gap-5">
            {updatedEvents.map((event, i) => (
              <Fragment key={event.id}>
                <EventAssignCard event={event} organizers={organizers ?? []} onAssign={handleAssign} />
                {i !== events.length - 1 && <Separator />}
              </Fragment>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
