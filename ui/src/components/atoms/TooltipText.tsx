import { ComponentProps } from "react";
import { Tooltip as TooltipPrimitive } from "radix-ui"
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";

interface Props extends ComponentProps<typeof TooltipPrimitive.Root> {
  text: string;
}

export function TooltipText({ text, children, ...props }: Props) {
  return (
    <Tooltip {...props}>
      <TooltipTrigger asChild>
        {children}
      </TooltipTrigger>
      <TooltipContent>
        {text}
      </TooltipContent>
    </Tooltip>
  )
}
