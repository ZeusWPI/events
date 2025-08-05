import { Input } from "../ui/input";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "../ui/dialog";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { ChangeEvent, useEffect, useState } from "react";
import { usePosterCreate, usePosterGetFile, usePosterUpdate } from "@/lib/api/poster";
import { Indeterminate } from "../atoms/Indeterminate";
import { Poster } from "@/lib/types/poster";
import { FileImg } from "../atoms/FileImg";
import { toast } from "sonner";

interface Props {
  poster: Poster;
  open: boolean;
  setOpen: (open: boolean) => void;
}

export function EventPosterDialog({ poster, open, setOpen }: Props) {
  const [submitting, setSubmitting] = useState(false)

  const { data: initialFile, isLoading } = usePosterGetFile(poster)
  const [file, setFile] = useState<File | undefined>(initialFile)

  useEffect(() => {
    if (open) {
      setFile(initialFile)
    } else {
      setFile(undefined)
    }
  }, [open, initialFile])

  const posterCreate = usePosterCreate()
  const posterUpdate = usePosterUpdate()

  const handleSubmit = () => {
    if (!file) {
      return
    }

    setSubmitting(true)

    let action
    let message
    if (poster.id) {
      action = posterUpdate
      message = "updated"
    } else {
      action = posterCreate
      message = "created"
    }

    action.mutate({ poster, file }, {
      onSuccess: () => {
        toast.success(`Poster ${message}`)
        setOpen(false)
      },
      onError: (error: Error) => toast.error("Failed", { description: error.message }),
      onSettled: () => setSubmitting(false)
    })
  }

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? undefined

    setFile(file)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="">
        <DialogTitle>
          Edit event poster
        </DialogTitle>
        <DialogDescription asChild>
          <div className="flex flex-col gap-8">
            <div className="space-y-3">
              <Label htmlFor="picture-scc">Poster</Label>
              {!isLoading ? (
                <>
                  <div className="flex items-center gap-1">
                    <Input id="picture-scc" type="file" accept=".png" onChange={handleChange} />
                  </div>
                  <div className="aspect-poster rounded-xl overflow-hidden">
                    {file && <FileImg file={file} isLoading={isLoading} alt="Poster preview" />}
                  </div>
                </>
              ) : (
                <Indeterminate />
              )}
            </div>
            <div className="flex justify-end">
              <Button onClick={handleSubmit} disabled={!file || submitting}>
                Submit
              </Button>
            </div>
          </div>
        </DialogDescription>
      </DialogContent>
    </Dialog>
  )
}
