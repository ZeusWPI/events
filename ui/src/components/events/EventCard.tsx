import { useAuth } from "@/lib/hooks/useAuth";
import { Check, CheckStatus, checkStatusToIcon } from "@/lib/types/check";
import type { Event } from "@/lib/types/event";
import { formatDate } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { ClipboardList, UserRound } from "lucide-react";
import { ComponentProps } from "react";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";
import { OrganizerGroup } from "../molecules/OrganizerGroup";

interface Props extends ComponentProps<'div'> {
  event: Event;
}

function checksColor(finishedChecks: Check[], checks: Check[]) {
  if (finishedChecks.length !== checks.length) {
    return "text-red-500"
  }

  return "text-secondary-foreground"
}

export function EventCard({ event, ...props }: Props) {
  const navigate = useNavigate();

  const { user } = useAuth();
  const isOrganizer = event.organizers.map(({ id }) => id).includes(user?.id ?? 0);

  const handleOnClick = () =>
    void navigate({ to: "/events/$id", params: { id: event.id.toString() } });

  const finishedChecks = event.checks.filter(check => check.status === CheckStatus.Done || check.status === CheckStatus.DoneLate)

  return (
    <Card onClick={handleOnClick} className={`transition-transform duration-300 hover:scale-102 hover:cursor-pointer ${isOrganizer && "border-primary"}`} {...props}>
      <CardHeader className="grow">
        <CardTitle>{event.name}</CardTitle>
        <CardDescription>
          <span>{formatDate(event.startTime)}</span>
        </CardDescription>
      </CardHeader>
      <CardFooter>
        <div className="flex items-center justify-between w-full">
          <Tooltip>
            <TooltipTrigger>
              <OrganizerGroup organizers={event.organizers} />
            </TooltipTrigger>
            {event.organizers.length > 0 && (
              <TooltipContent>
                <div className="flex flex-col">
                  {event.organizers.map(organizer => (
                    <div key={organizer.id} className="flex items-center gap-1">
                      <UserRound className="size-3" />
                      <span>{organizer.name}</span>
                    </div>
                  ))}
                </div>
              </TooltipContent>
            )}
          </Tooltip>
          <Tooltip>
            <TooltipTrigger>
              <div className="flex items-center space-x-2">
                <span className={checksColor(finishedChecks, event.checks)}>{`${finishedChecks.length}/${event.checks.length}`}</span>
                <ClipboardList className="size-6" />
              </div>
            </TooltipTrigger>
            {event.checks.length > 0 && (
              <TooltipContent>
                <div className="flex flex-col">
                  {event.checks.map(check => (
                    <div key={check.id} className="flex items-center gap-1">
                      <div className="size-3 flex items-center">
                        {checkStatusToIcon[check.status]}
                      </div>
                      <span>{check.description}</span>
                    </div>
                  ))}
                </div>
              </TooltipContent>
            )}
          </Tooltip>
        </div>
      </CardFooter>
    </Card>
  );
}
