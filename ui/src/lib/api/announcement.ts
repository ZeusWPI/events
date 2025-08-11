import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Announcement, AnnouncementSchema, convertAnnouncementsToModel } from "../types/announcement";
import { Year } from "../types/year";
import { apiDelete, apiGet, apiPost, apiPut } from "./query";

const ENDPOINT = "announcement";
const STALE_5_MIN = 5 * 60 * 1000;

export function useAnnouncementByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["announcement", id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/year/${id}`, convertAnnouncementsToModel)).data,
    staleTime: STALE_5_MIN,
    throwOnError: true,
    enabled: id > 0,
  })
}

export function useAnnouncementCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (announcement: AnnouncementSchema) => apiPut(ENDPOINT, announcement),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["announcement"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    }
  })
}

export function useAnnouncementUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (announcement: AnnouncementSchema) => apiPost(`${ENDPOINT}/${announcement.id}`, announcement),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["announcement"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    }
  })
}

export function useAnnouncementDelete() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id }: Pick<Announcement, "id">) => apiDelete(`${ENDPOINT}/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["announcement"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    }
  })
}
