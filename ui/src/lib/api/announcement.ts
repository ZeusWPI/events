import { useMutation, useQueryClient } from "@tanstack/react-query";
import { convertAnnouncementToJSON, Announcement } from "../types/announcement";
import { apiPost, apiPut } from "./query";

const ENDPOINT = "announcement";

export function useAnnouncementCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (announcement: Announcement) => apiPut(ENDPOINT, convertAnnouncementToJSON(announcement)),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}

export function useAnnouncementUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (announcement: Announcement) => apiPost(`${ENDPOINT}/${announcement.id}`, convertAnnouncementToJSON(announcement)),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}
