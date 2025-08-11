import { Indeterminate } from "@/components/atoms/Indeterminate";
import { useYear, useYearLock } from "@/lib/hooks/useYear";
import { useNavigate, useParams } from "@tanstack/react-router";
import Error404 from "../404";
import { AnnouncementForm } from "@/components/announcements/AnnouncementForm";
import { AnnouncementSchema } from "@/lib/types/announcement";
import { useAnnouncementByYear, useAnnouncementUpdate } from "@/lib/api/announcement";
import { toast } from "sonner";

export function AnnouncementsEdit() {
  const { announcementId } = useParams({ from: "/announcements/edit/$announcementId" })

  const { year } = useYear()
  const { data: announcements, isLoading: isLoadingAnnouncements } = useAnnouncementByYear(year)
  const announcement = announcements?.find(a => a.id === Number(announcementId))

  useYearLock()
  const navigate = useNavigate()
  const update = useAnnouncementUpdate()

  if (isLoadingAnnouncements) {
    return <Indeterminate />
  }

  if (!announcement) {
    return <Error404 />
  }

  const handleSubmit = (announcement: AnnouncementSchema) => {
    update.mutate(announcement, {
      onSuccess: () => {
        toast.success("Success")
        navigate({ to: "/announcements" })
      },
      onError: error => toast.error("Failed", { description: error.message }),
    })
  }

  return (
    <AnnouncementForm announcement={announcement} onSubmit={handleSubmit} />
  )
}
