import type { Event } from "@/lib/types/event";
import type { Organizer } from "@/lib/types/organizer";
import { formatDate } from "@/lib/utils/utils";
import { CircleAlert } from "lucide-react";
import { useState } from "react";
import { MultiSelect } from "../organisms/MultiSelect";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";

interface Props {
  event: Event;
  organizers: Organizer[];
  onAssign: (eventId: number, selectedOrganizers: number[]) => void;
}

export function EventAssignCard({ event, organizers, onAssign }: Props) {
  const [selected, setSelected] = useState(event.organizers.map(({ id }) => id.toString()));

  const onValueChange = (value: string[]) => {
    setSelected(value);
    onAssign(event.id, value.map(Number));
  };

  return (
    <div className="grid lg:grid-cols-3 gap-2">
      <div className="flex flex-col">
        <div className="flex items-center gap-2">
          {selected.length === 0 && (
            <Tooltip>
              <TooltipTrigger>
                <CircleAlert color="red" />
              </TooltipTrigger>
              <TooltipContent>
                No organizers assigned
              </TooltipContent>
            </Tooltip>
          )}
          <span className="text-md break-words overflow-hidden">{event.name}</span>
        </div>
        <span className="text-sm/6 text-muted-foreground">{formatDate(event.startTime)}</span>
      </div>
      <MultiSelect
        options={organizers.map(({ id, name }) => ({ value: id.toString(), label: name }))}
        onValueChange={onValueChange}
        defaultValue={selected}
        placeholder="Select organizers"
        className="lg:col-span-2"
      />
    </div>
  );
}
