import { Announcement } from "@/lib/types/announcement";
import { AnnouncementCard } from "./AnnouncementCard";
import { ComponentProps } from "react";

interface Props extends ComponentProps<'div'> {
  announcements: Announcement[];
}

export function AnnouncementList({ announcements, ...props }: Props) {
  return (
    <div className="grid gap-4 grid-cols-1" {...props}>
      {announcements.map(a => <AnnouncementCard key={a.id} announcement={a} />)}
    </div>
  )
}
