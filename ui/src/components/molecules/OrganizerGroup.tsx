import { Organizer } from "@/lib/types/organizer";
import { OrganizerIcon } from "../atoms/OrganizerIcon";

interface Props {
  organizers: Organizer[];
}

export function OrganizerGroup({ organizers }: Props) {
  if (!organizers.length) {
    return <div className="text-red-500">No users</div>
  }

  return (
    <div className="*:data-[slot=avatar]:ring-background flex -space-x-2 *:data-[slot=avatar]:ring-2">
      {organizers.map(o => <OrganizerIcon key={o.id} user={o} />)}
    </div>
  )
}
