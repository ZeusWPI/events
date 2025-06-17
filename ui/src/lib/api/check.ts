import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Check } from "../types/check";
import { apiDelete, apiPost } from "./query";

const ENDPOINT = "check";

export function useCheckCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (check: Pick<Check, 'eventId' | 'description'>) => apiPost(`${ENDPOINT}`, { event_id: check.eventId, description: check.description }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] }),
  })
}

export function useCheckToggle() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (check: Pick<Check, 'id'>) => apiPost(`${ENDPOINT}/${check.id}`),
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
