import { usePosterDelete, usePosterGetFile } from "@/lib/api/poster";
import { Poster } from "@/lib/types/poster";
import { Year } from "@/lib/types/year";
import { DownloadIcon, PencilIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { FileImg } from "../atoms/FileImg";
import { TooltipText } from "../atoms/TooltipText";
import { DeleteConfirm } from "../molecules/DeleteConfirm";
import { Button } from "../ui/button";
import { EventPosterDialog } from "./EventPosterDialog";
import { getUuid } from "@/lib/utils/utils";

interface Props {
  title: string;
  description: string;
  poster: Poster;
  year: Year;
}

export const EventPoster = ({ title, description, poster, year }: Props) => {
  const { data: file, isLoading } = usePosterGetFile(poster)

  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)

  const posterDelete = usePosterDelete()

  const handleDownload = () => {
    const a = document.createElement('a')

    a.href = `/api/poster/${poster.id}/file?original=true`
    a.download = `${getUuid()}.png`
    document.body.appendChild(a)
    a.click()

    document.body.removeChild(a)
  }

  const handleDeleteConfirm = () => {
    posterDelete.mutate({ poster, year }, {
      onSuccess: () => toast.success("Poster deleted"),
      onError: (error: Error) => toast.error("Failed", { description: error.message }),
      onSettled: () => setOpenDelete(false)
    })
  }

  return (
    <div className="flex flex-col border rounded-xl max-w-96 w-full justify-between">
      <div>
        {poster.id > 0 && (
          <div className="rounded-xl overflow-hidden aspect-poster">
            <FileImg file={file} isLoading={isLoading} alt={title} />
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
          {file && (
            <TooltipText text="Download uncompressed poster">
              <Button onClick={handleDownload} size="icon" variant="ghost" className="size-6">
                <DownloadIcon />
              </Button>
            </TooltipText>
          )}
          <Button onClick={() => setOpenEdit(true)} size="icon" variant="ghost" disabled={isLoading} className="size-6">
            {file
              ? <PencilIcon />
              : <PlusIcon />
            }
          </Button>
          <Button onClick={() => setOpenDelete(true)} size="icon" variant="secondary" disabled={!file || isLoading} className="size-6">
            <Trash2Icon className="text-red-500" />
          </Button>
        </div>
      </div>
      <EventPosterDialog
        poster={poster}
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
