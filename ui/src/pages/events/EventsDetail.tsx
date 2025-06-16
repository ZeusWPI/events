import { useParams } from "@tanstack/react-router";
import { Link, PlusIcon, UserRound } from "lucide-react";
import { Title } from "@/components/atoms/Title";
import { Datalist, DatalistItem, DatalistItemContent, DatalistItemTitle } from "@/components/molecules/Datalist";
import { HeadlessCard } from "@/components/molecules/HeadlessCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useEventByYear } from "@/lib/api/event";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { formatDate } from "@/lib/utils/utils";
import Error404 from "../404";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { CheckCard } from "@/components/check/CheckCard";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { useCheckCreate } from "@/lib/api/check";
import { toast } from "sonner";

export function EventsDetail() {
  const { year: yearString, id: eventID } = useParams({ from: "/events/$year/$id" });

  const { data: years } = useYearGetAll();
  // Event component makes sure it exists
  const year = years!.find(({ formatted }) => formatted === yearString)!;
  const { data: events, isLoading } = useEventByYear(year);
  const event = events?.find(event => event.id.toString() === eventID);

  const [adding, setAdding] = useState(false)
  const [addCheck, setAddCheck] = useState(false)
  const [addCheckDescription, setAddCheckDescription] = useState("")

  const checkCreate = useCheckCreate()

  const cancelCheck = () => {
    setAddCheck(false)
    setAddCheckDescription("")
  }
  const createCheck = () => {
    if (addCheckDescription === "") {
      return
    }

    setAdding(true)

    checkCreate.mutate({ eventId: Number(eventID), description: addCheckDescription }, {
      onSuccess: () => {
        toast.success("Success")
        setAddCheck(false)
        setAddCheckDescription("")
      },
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setAdding(false)
    })
  }

  useBreadcrumb({ title: event?.name ?? "", link: { to: "/events/$year/$id", params: { year: yearString, id: eventID } } });

  if (isLoading) {
    return <Indeterminate />
  }

  if (!event) {
    return <Error404 />;
  }

  const finishedChecks = event.checks.filter(check => check.done)

  return (
    <div className="grid lg:grid-cols-3 gap-8">
      <PageHeader className="col-span-full">
        <Title>{event.name}</Title>
        <Button variant="outline" size="icon" asChild>
          <a href={event.url} rel="noopener noreferrer" target="_blank">
            <Link />
          </a>
        </Button>
      </PageHeader>
      <div className="lg:col-span-2">
        <HeadlessCard>
          <CardHeader className="px-4 sm:px-0 pt-0">
            <CardTitle>General</CardTitle>
          </CardHeader>
          <CardContent className="px-0">
            <Datalist>
              <DatalistItem>
                <DatalistItemTitle>Description</DatalistItemTitle>
                <DatalistItemContent>{event.description}</DatalistItemContent>
              </DatalistItem>
              <DatalistItem>
                <DatalistItemTitle>Time</DatalistItemTitle>
                <DatalistItemContent>
                  <div className="flex flex-col">
                    <span>{formatDate(event.startTime)}</span>
                    {event.endTime && (
                      <span>{formatDate(event.endTime)}</span>
                    )}
                  </div>
                </DatalistItemContent>
              </DatalistItem>
              <DatalistItem>
                <DatalistItemTitle>Location</DatalistItemTitle>
                <DatalistItemContent>{event.location}</DatalistItemContent>
              </DatalistItem>
            </Datalist>
          </CardContent>
        </HeadlessCard>
      </div>
      <div className="lg:col-span-1">
        <div className="flex h-full">
          <Separator orientation="vertical" className="h-full hidden lg:block" />
          <HeadlessCard className="h-full mx-auto">
            <CardHeader className="px-0 pt-0">
              <CardTitle className="flex gap-2 items-center">
                <UserRound />
                Organizers
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col space-y-1 px-0">
              {event.organizers.length
                ? event.organizers.map(organizer => (
                  <span key={organizer.id}>{organizer.name}</span>
                ))
                : (
                  <span>No one assigned</span>
                )}
            </CardContent>
          </HeadlessCard>
        </div>
      </div>
      <div className="col-span-full">
        <HeadlessCard>
          <CardHeader className="px-4 sm:px-0 pt-0">
            <CardTitle>
              <div className="flex justify-between">
                <span>
                  Checks <span className="font-normal">{`${finishedChecks.length}/${event.checks.length}`}</span>
                </span>
                {!addCheck && (
                  <Button onClick={() => setAddCheck(true)} size="icon" variant="ghost">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <PlusIcon className="text-primary" />
                      </TooltipTrigger>
                      <TooltipContent>
                        <span>Add a new check</span>
                      </TooltipContent>
                    </Tooltip>
                  </Button>
                )}
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent className="px-0">
            <div className="flex flex-col gap-8">
              {addCheck && (
                <div className="flex w-full gap-2">
                  <Input placeholder="Description" onChange={e => setAddCheckDescription(e.target.value)} />
                  <Button onClick={cancelCheck} variant="outline">
                    <span>Cancel</span>
                  </Button>
                  <Button onClick={createCheck} disabled={adding}>
                    <span>Create</span>
                  </Button>
                </div>
              )}
              {event.checks.map(check => (
                <CheckCard key={check.id} check={check} />
              ))}
            </div>
          </CardContent>
        </HeadlessCard>
      </div>
    </div>
  );
}
