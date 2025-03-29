import { cn } from "@/lib/utils/utils";
import { Badge } from "../ui/badge";

type Color = "red" | "green";

const colors: Record<Color, string> = {
  red: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  green: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
};

interface Props {
  color: Color;
}

export function Pill({ color, ...props }: Props & React.ComponentProps<"div">) {
  return (
    <Badge variant="outline" className={cn("text-xs gap-1.5 rounded-full", colors[color])} {...props} />
  );
}
