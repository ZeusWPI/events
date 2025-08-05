import { Button } from "../ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Dialog as DialogPrimitive } from "radix-ui"

interface Props extends React.ComponentProps<typeof DialogPrimitive.Root> {
  onDelete: () => void;
}

export function DeleteConfirm({ onDelete, open, onOpenChange }: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            Delete confirmation
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete it?
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button onClick={onDelete} variant="destructive">Delete</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
