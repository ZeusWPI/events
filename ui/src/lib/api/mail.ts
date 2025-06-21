import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { apiGet, apiPost, apiPut } from "./query"
import { convertMailsToModel, Mail } from "../types/mail"

const ENDPOINT = "mail"
const STALE_5_MIN = 5 * 60 * 1000

export function useMailGetAll() {
  return useQuery({
    queryKey: ["mail"],
    queryFn: async () => (await apiGet(ENDPOINT, convertMailsToModel)).data,
    staleTime: STALE_5_MIN,
    throwOnError: true,
  })
}

export function useMailCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (args: { mail: Pick<Mail, 'title' | 'content' | 'sendTime'>, eventIds: number[] }) => apiPut(ENDPOINT, { ...args.mail, eventIds: args.eventIds }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mail"] })
      queryClient.invalidateQueries({ queryKey: ["event"] })
    },
  })
}

export function useMailUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (args: { mail: Pick<Mail, 'id' | 'title' | 'content' | 'sendTime'>, eventIds: number[] }) => apiPost(`${ENDPOINT}/${args.mail.id}`, { ...args.mail, eventIds: args.eventIds }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mail"] })
      queryClient.invalidateQueries({ queryKey: ["event"] })
    },
  })
}
