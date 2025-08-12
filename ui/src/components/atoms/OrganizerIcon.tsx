import { Organizer } from "@/lib/types/organizer";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { cn } from "@/lib/utils/utils";

interface Props {
  user: Organizer;
  className?: string;
}

export function OrganizerIcon({ user, className }: Props) {
  const initials = user?.name.split(" ").map(n => n[0]).join("");

  return (
    <Avatar className={cn("aspect-square h-8 w-8 rounded-lg", className)}>
      <AvatarImage src={`https://zpi.zeus.gent/image/${user.zauthId}?size=64&placeholder=true`} alt={initials} />
      <AvatarFallback className="rounded-lg">{initials}</AvatarFallback>
    </Avatar>
  )
}
