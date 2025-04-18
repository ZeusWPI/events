import { cn } from "@/lib/utils/utils";

export function DividerText({ className, ...props }: React.ComponentProps<"span">) {
  return (
    <div className={cn("flex py-4 items-center", className)}>
      <div className="grow border-t" />
      <span className={cn("shrink mx-4 text-muted-foreground", className)} {...props}>
      </span>
      <div className="grow border-t" />
    </div>
  );
}
