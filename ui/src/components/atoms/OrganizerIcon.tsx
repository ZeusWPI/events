import { Organizer } from "@/lib/types/organizer";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { cn } from "@/lib/utils/utils";
import { TooltipText } from "./TooltipText";
import { useTheme } from "@/lib/hooks/useTheme";

interface Props {
  user: Organizer;
  tooltip?: boolean;
  className?: string;
}

export function OrganizerIcon({ user, tooltip = false, className }: Props) {
  const { actualTheme } = useTheme()
  const background = actualTheme === "light" ? "bg-white" : "bg-black"

  const initials = user?.name.split(" ").map(n => n[0]).join("");

  const Icon = (
    <Avatar className={cn(`aspect-square h-8 w-8 rounded-lg bg-inherit ${background}`, className)}>
      <AvatarImage src={`https://zpi.zeus.gent/image/${user.zauthId}?size=64&placeholder=true`} alt={initials} />
      <AvatarFallback className="rounded-lg">{initials}</AvatarFallback>
    </Avatar>
  )

  if (!tooltip) return Icon

  return (
    <TooltipText text={user.name}>
      {Icon}
    </TooltipText>
  )
}
