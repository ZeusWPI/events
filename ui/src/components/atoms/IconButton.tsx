import { ComponentProps } from "react";
import { Button } from "../ui/button";

interface Props extends ComponentProps<'button'> {
  asChild?: boolean;
}

export function IconButton({ ...props }: Props) {
  return <Button size="icon" variant="outline" {...props} />
}
