import { AnyFieldApi } from "@tanstack/react-form";
import { ComponentProps } from "react";

interface Props extends ComponentProps<'div'> {
  field: AnyFieldApi;
}

export function FormField({ field, children, ...props }: Props) {
  return (
    <div {...props}>
      {children}
      {field.state.meta.isTouched && !field.state.meta.isValid && field.state.meta.errors.length && (
        <em className="text-red-500">{field.state.meta.errors[0].message}</em>
      )}
    </div>
  )
}
