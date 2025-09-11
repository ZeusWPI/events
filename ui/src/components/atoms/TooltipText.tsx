import { ComponentProps } from "react";
import { Tooltip as TooltipPrimitive } from "radix-ui"
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";
import { Separator } from "../ui/separator";

interface Props extends ComponentProps<typeof TooltipPrimitive.Root> {
  text: string;
  subtext?: string;
}

export function TooltipText({ text, subtext, children, ...props }: Props) {
  return (
    <Tooltip {...props}>
      <TooltipTrigger asChild>
        {children}
      </TooltipTrigger>
      <TooltipContent>
        <div className="space-y-2">
          <p className="whitespace-pre">
            {text}
          </p>
          {subtext && (
            <>
              <Separator />
              <p className="whitespace-pre text-muted-foreground">
                {subtext}
              </p>
            </>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  )
}
