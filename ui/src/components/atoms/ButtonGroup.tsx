import { cn } from "@/lib/utils/utils";
import { ComponentProps } from "react";

export function ButtonGroup({ className, ...props }: ComponentProps<'div'>) {
  return <div className={cn("flex gap-2 items-center", className)} {...props} />
}
