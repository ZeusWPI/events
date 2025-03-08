import { cn } from "@/lib/utils/utils";
import { Card } from "../ui/card";

export function HeadlessCard({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <Card className={cn("shadow-none border-none", className)} {...props} />
  );
}
