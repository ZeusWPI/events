import { PencilIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { FileImg } from "../atoms/FileImg";
import { Button } from "../ui/button";
import { Poster } from "@/lib/types/poster";
import { usePosterDelete, usePosterGetFile } from "@/lib/api/poster";
import { useState } from "react";
import { EventPosterDialog } from "./EventPosterDialog";
import { DeleteConfirm } from "../molecules/DeleteConfirm";
import { toast } from "sonner";
import { Year } from "@/lib/types/year";

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

  const handleDeleteConfirm = () => {
    posterDelete.mutate({ poster, year }, {
      onSuccess: () => toast.success("Poster deleted"),
      onError: (error: Error) => toast.error("Failed", { description: error.message }),
      onSettled: () => setOpenDelete(false)
    })
  }

  return (
    <div className="flex flex-col border rounded-xl max-w-96 w-full">
      {poster.id > 0 && (
        <div className="rounded-xl overflow-hidden aspect-poster">
          <FileImg file={file} isLoading={isLoading} alt={title} />
        </div>
      )}
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
