import { Event } from "@/lib/types/event";
import { XIcon } from "lucide-react";
import { Input } from "../ui/input";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "../ui/dialog";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { ChangeEvent, useEffect, useState } from "react";
import { usePosterCreate, usePosterDelete, usePosterGetFile, usePosterUpdate } from "@/lib/api/poster";
import { Indeterminate } from "../atoms/Indeterminate";
import { toast } from "sonner";
import { arrayEqual } from "@/lib/utils/utils";

interface Props {
  event: Event;
  open: boolean;
  setOpen: (open: boolean) => void;
}

export function EventPosterDialog({ event, open, setOpen }: Props) {
  const [submitting, setSubmitting] = useState(false)

  const { data: oldBig, isLoading: isLoadingBig } = usePosterGetFile(event.posters.find(p => !p.scc)?.id ?? 0, event.id)
  const { data: oldScc, isLoading: isLoadingScc } = usePosterGetFile(event.posters.find(p => p.scc)?.id ?? 0, event.id)

  const [big, setBig] = useState<File | null>(null)
  const [scc, setScc] = useState<File | null>(null)

  const bigURL = big ? URL.createObjectURL(big) : ""
  const sccURL = scc ? URL.createObjectURL(scc) : ""

  const posterCreate = usePosterCreate()
  const posterUpdate = usePosterUpdate()
  const posterDelete = usePosterDelete()

  useEffect(() => {
    if (oldBig) {
      setBig(oldBig)
    }

    if (oldScc) {
      setScc(oldScc)
    }
  }, [oldBig, oldScc])

  const handleChange = (e: ChangeEvent<HTMLInputElement>, scc: boolean) => {
    const file = e.target.files?.[0] ?? null;
    if (scc) {
      setScc(file);
    } else {
      setBig(file);
    }
  }

  const handleSubmit = async () => {
    setSubmitting(true);

    const posterTypes = [
      { key: "Big", file: big, oldFile: oldBig, scc: false },
      { key: "Scc", file: scc, oldFile: oldScc, scc: true },
    ];

    let requestSend = false

    for (let i = 0; i < posterTypes.length; i++) {
      const { key, file, oldFile, scc } = posterTypes[i]!;

      const data = (await file?.bytes()) ?? new Uint8Array();
      const oldData = (await oldFile?.bytes()) ?? new Uint8Array();

      if (arrayEqual(data, oldData)) {
        continue
      }

      requestSend = true

      const poster = {
        id: event.posters.find((p) => p.scc === scc)?.id ?? 0,
        eventId: event.id,
        scc,
      };

      const commonMutationOptions = {
        onSuccess: () => toast.success(`${key} ${file ? (oldFile ? "updated" : "created") : "deleted"}`),
        onError: (error: Error) => toast.error("Failed", { description: error.message }),
        onSettled: () => setSubmitting(false),
      };

      if (file) {
        const mutation = oldFile ? posterUpdate : posterCreate;
        mutation.mutate({ poster, file }, commonMutationOptions);
      } else {
        posterDelete.mutate(poster, commonMutationOptions);
      }
    }

    if (!requestSend) {
      setSubmitting(false)
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="lg:min-w-4xl">
        <DialogTitle>
          Edit event posters
        </DialogTitle>
        <DialogDescription asChild>
          <div className="flex flex-col gap-8">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
              <div className="space-y-3">
                <Label htmlFor="picture-big">Big</Label>
                {!isLoadingBig ? (
                  <>
                    <div className="flex items-center gap-1">
                      <Input id="picture-big" type="file" accept=".png" onChange={e => handleChange(e, false)} />
                      <Button size="icon" onClick={() => setBig(null)} variant="ghost">
                        <XIcon />
                      </Button>
                    </div>
                    {big && <img src={bigURL} alt="Big preview" />}
                  </>
                ) : (
                  <Indeterminate />
                )}
              </div>
              <div className="space-y-3">
                <Label htmlFor="picture-scc">Scc</Label>
                {!isLoadingScc ? (
                  <>
                    <div className="flex items-center gap-1">
                      <Input id="picture-scc" type="file" accept=".png" onChange={e => handleChange(e, true)} />
                      <Button size="icon" onClick={() => setScc(null)} variant="ghost">
                        <XIcon />
                      </Button>
                    </div>
                    {scc && <img src={sccURL} alt="Scc preview" />}
                  </>
                ) : (
                  <Indeterminate />
                )}
              </div>
            </div>
            <div className="flex justify-end">
              <Button onClick={handleSubmit} disabled={submitting}>
                Submit
              </Button>
            </div>
          </div>
        </DialogDescription>
      </DialogContent>
    </Dialog>
  )
}
