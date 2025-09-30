import { usePosterDelete } from "@/lib/api/poster";
import { Poster } from "@/lib/types/poster";
import { Year } from "@/lib/types/year";
import { getUuid } from "@/lib/utils/utils";
import { DownloadIcon, PencilIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { LoadableImage } from "../atoms/LoadableImage";
import { TooltipText } from "../atoms/TooltipText";
import { DeleteConfirm } from "../molecules/DeleteConfirm";
import { Button } from "../ui/button";
import { EventPosterDialog } from "./EventPosterDialog";

interface Props {
  title: string;
  description: string;
  poster?: Poster;
  eventId: number;
  scc: boolean;
  year: Year;
}

export const EventPoster = ({ title, description, poster, eventId, scc, year }: Props) => {
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)

  const posterDelete = usePosterDelete()

  const handleDownload = () => {
    if (!poster) return

    const a = document.createElement('a')

    a.href = `/api/poster/${poster.id}?original=true`
    a.download = `${getUuid()}.png`
    document.body.appendChild(a)
    a.click()

    document.body.removeChild(a)
  }

  const handleDeleteConfirm = () => {
    if (!poster) return

    posterDelete.mutate({ poster, year }, {
      onSuccess: () => toast.success("Poster deleted"),
      onError: (error: Error) => toast.error("Failed", { description: error.message }),
      onSettled: () => setOpenDelete(false)
    })
  }

  return (
    <div className="flex flex-col border rounded-xl max-w-96 w-full justify-between mx-auto lg:mx-0">
      <div>
        {poster && (
          <div className="rounded-xl overflow-hidden aspect-poster">
            <LoadableImage src={`/api/poster/${poster.id}`} alt={title} />
          </div>
        )}
      </div>
      <div className="flex justify-between p-4">
        <div className="flex flex-col space-y-1">
          <span className="leading-none font-semibold">
            {title}
          </span>
          <span className="text-muted-foreground text-sm">
            {description}
          </span>
        </div>
        <div className="flex gap-1 items-end">
          {poster && (
            <TooltipText text="Download uncompressed poster">
              <Button onClick={handleDownload} size="icon" variant="ghost" className="size-6">
                <DownloadIcon />
              </Button>
            </TooltipText>
          )}
          <TooltipText text={poster ? "Edit" : "Create"}>
            <Button onClick={() => setOpenEdit(true)} size="icon" variant="ghost" className="size-6">
              {poster
                ? <PencilIcon />
                : <PlusIcon />
              }
            </Button>
          </TooltipText>
          <TooltipText text="Delete poster">
            <Button onClick={() => setOpenDelete(true)} size="icon" variant="secondary" disabled={!poster} className="size-6">
              <Trash2Icon className="text-red-500" />
            </Button>
          </TooltipText>
        </div>
      </div>
      <EventPosterDialog
        poster={poster ?? { id: 0, eventId, scc }}
        open={openEdit}
        setOpen={setOpenEdit}
      />
      <DeleteConfirm
        open={openDelete}
        onOpenChange={setOpenDelete}
        onDelete={handleDeleteConfirm}
      />
    </div>
  )
}
