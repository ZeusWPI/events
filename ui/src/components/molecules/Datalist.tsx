import { cn } from "@/lib/utils/utils";

export function Datalist({ className, ...props }: React.ComponentProps<"dl">) {
  return (
    <dl
      className={cn("divide-y divide-muted-foreground border-t border-muted-foreground", className)}
      {...props}
    />
  );
}

export function DatalistItem({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("px-4 py-6 sm:grid sm:grid-cols-3 sm-gap-4 sm:px-0", className)}
      {...props}
    />
  );
}

export function DatalistItemTitle({ className, ...props }: React.ComponentProps<"dt">) {
  return (
    <dt
      className={cn("font-medium", className)}
      {...props}
    />
  );
}

export function DatalistItemContent({ className, ...props }: React.ComponentProps<"dd">) {
  return (
    <dd
      className={cn("sm:col-span-2 sm:mt-0 text-muted-foreground", className)}
      {...props}
    />
  );
}
