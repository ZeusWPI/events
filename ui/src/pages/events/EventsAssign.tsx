import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { EventAssignCard } from "@/components/events/EventAssignCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useEventByYear, useEventSaveOrganizers } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useYearGetAll } from "@/lib/api/year";
import { useIsMobile } from "@/lib/hooks/use-mobile";
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

  const { mutate: updateOrganizers } = useEventSaveOrganizers();

  const [updatedEvents, setUpdatedEvents] = useState(events ?? []);
  const [isDirty, setIsDirty] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  useBreadcrumb({ title: "Assign", link: { to: "/events/$year/assign", params: { year: yearString } } });
  const isMobile = useIsMobile();

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
    updateOrganizers(updatedEvents, {
      onSuccess: () => {
        setIsDirty(false);
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
    setIsDirty(true);
  };

  if (!events) {
    return <Indeterminate />;
  }

  return (
    <div className="grid xl:grid-cols-4 gap-8">
      <PageHeader className="col-span-full">
        <Title>{`Assign${!isMobile ? ` to Events ${yearString}` : ""}`}</Title>
        <div className="flex items-center gap-6">
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
      </PageHeader>
      <div className="sticky top-6">
        <Card className={`w-full xl:max-w-80 sticky top-6 ${isMobile ? "max-h-48 overflow-y-scroll" : ""}`}>
          <CardContent className="flex flex-col gap-1">
            {organizersEventCount.map((organizer, index) => (
              <motion.div key={organizer.id} layout className="flex justify-between items-center w-full">
                <span className={index % 2 === 1 ? "text-muted-foreground" : ""}>{organizer.name}</span>
                <span className={index % 2 === 1 ? "text-muted-foreground" : ""}>{organizer.events}</span>
              </motion.div>
            ))}
          </CardContent>
        </Card>
      </div>
      <div className="xl:col-span-3 flex flex-col gap-5">
        {updatedEvents.map((event, i) => (
          <Fragment key={event.id}>
            <EventAssignCard event={event} organizers={organizers ?? []} onAssign={handleAssign} />
            {i !== events.length - 1 && <Separator />}
          </Fragment>
        ))}
      </div>
    </div>
  );
}
