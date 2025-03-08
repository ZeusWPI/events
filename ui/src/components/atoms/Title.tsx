import { cn } from "@/lib/utils/utils";

export function Title({ className, ...props }: React.ComponentProps<"h1">) {
  return (
    <h1 className={cn("text-3xl font-bold leading-none tracking-tight", className)} {...props} />
  );
}
