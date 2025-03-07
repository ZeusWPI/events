import type { Event } from "@/lib/types/types";
import { formatDate } from "@/lib/utils/utils";
import { ClipboardList, UserRound } from "lucide-react";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";

interface Props {
  event: Event;
}

function organizersColor(amount: number) {
  switch (amount) {
    case 0:
      return "text-red-500";
    case 1:
      return "text-orange-500";
    default:
      return "text-green-500";
  }
}

export function EventCard({ event }: Props) {
  return (
    <Card className="transition-transform duration-300 hover:scale-105 hover:cursor-pointer">
      <CardHeader className="grow">
        <CardTitle>{event.name}</CardTitle>
        <CardDescription>
          <span>{formatDate(event.startTime)}</span>
        </CardDescription>
      </CardHeader>
      <CardFooter>
        <div className="flex items-center justify-between w-full">
          <div className="flex space-x-2">
            <UserRound className="size-6" />
            <span className={organizersColor(event.organizers.length)}>{event.organizers.length}</span>
          </div>
          <div className="flex items-center space-x-2">
            <span>0/5</span>
            <ClipboardList className="size-6" />
          </div>
        </div>
      </CardFooter>
    </Card>
  );
}
