import { Button } from "../ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Dialog as DialogPrimitive } from "radix-ui"

interface Props extends React.ComponentProps<typeof DialogPrimitive.Root> {
  onResend: () => void;
}

export function ResendConfirm({ onResend, open, onOpenChange }: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            Resend confirmation
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to resend it?
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button onClick={onResend} variant="destructive">Resend</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
