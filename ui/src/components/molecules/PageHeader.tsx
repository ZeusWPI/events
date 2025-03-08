import { cn } from "@/lib/utils/utils";
import React from "react";

export function PageHeader({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex justify-between items-end border-b border-primary pb-2", className)}
      {...props}
    />
  );
}
