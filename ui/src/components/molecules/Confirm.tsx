import { Button } from "../ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Dialog as DialogPrimitive } from "radix-ui"

interface Props extends React.ComponentProps<typeof DialogPrimitive.Root> {
  title: string;
  description: string;
  confirmText: string;
  onConfirm: () => void;
}

export function Confirm({ title, description, confirmText, onConfirm, open, onOpenChange }: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button onClick={onConfirm} variant="destructive">{confirmText}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
