import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Check, convertChecksToModel } from "../types/check";
import { Year } from "../types/year";
import { apiDelete, apiGet, apiPost, apiPut } from "./query";

const ENDPOINT = "check";
const STALE_5_MIN = 5 * 60 * 1000;

export function useCheckGetByYear(year: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["check", "year", year.id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/year/${year.id}`, convertChecksToModel)).data,
    staleTime: STALE_5_MIN,
    throwOnError: true,
    enabled: year.id > 0,
  })
}

export function useCheckCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (check: Pick<Check, 'eventId' | 'description'>) => apiPut(`${ENDPOINT}`, { event_id: check.eventId, description: check.description }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] }),
  })
}

export function useCheckUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (check: Pick<Check, "id" | "status" | "description">) => (await apiPost(`${ENDPOINT}/${check.id}`, check)).data,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}

export function useCheckDelete() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (check: Pick<Check, 'id'>) => apiDelete(`${ENDPOINT}/${check.id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}
