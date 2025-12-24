import { AnnouncementForm } from "@/components/announcements/AnnouncementForm";
import { useAnnouncementCreate } from "@/lib/api/announcement";
import { useYearLock } from "@/lib/hooks/useYear";
import { AnnouncementSchema } from "@/lib/types/announcement";
import { getRouteApi, useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";

const route = getRouteApi('/announcements/create')

export function AnnouncementsCreate() {
  const { eventIds } = route.useSearch()

  const navigate = useNavigate()
  const create = useAnnouncementCreate()

  useYearLock()

  const handleSubmit = (announcement: AnnouncementSchema) => {
    create.mutate(announcement, {
      onSuccess: () => {
        toast.success("Success")
        navigate({ to: "/announcements" })
      },
      onError: error => toast.error("Failed", { description: error.message }),
    })

  }

  return (
    <AnnouncementForm defaultEvents={eventIds} onSubmit={handleSubmit} />
  )
}
