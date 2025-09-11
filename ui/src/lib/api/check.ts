import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Check } from "../types/check";
import { apiDelete, apiPost, apiPut } from "./query";

const ENDPOINT = "check";

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
