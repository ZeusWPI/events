import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Announcement } from "../types/announcement";
import { apiPost, apiPut } from "./query";

const ENDPOINT = "announcement";

export function useAnnouncementCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (announcement: Announcement) => apiPut(ENDPOINT, announcement),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}

export function useAnnouncementUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (announcement: Announcement) => apiPost(`${ENDPOINT}/${announcement.id}`, announcement),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}
