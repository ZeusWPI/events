import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { convertMailsToModel, Mail, MailSchema } from "../types/mail"
import { Year } from "../types/year"
import { apiDelete, apiGet, apiPost, apiPut } from "./query"

const ENDPOINT = "mail"
const STALE_5_MIN = 5 * 60 * 1000

export function useMailByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["mail", id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/year/${id}`, convertMailsToModel)).data,
    staleTime: STALE_5_MIN,
    throwOnError: true,
    enabled: id > 0,
  })
}

export function useMailCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (mail: MailSchema) => apiPut(ENDPOINT, mail),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mail"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    }
  })
}

export function useMailUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (mail: MailSchema) => apiPost(`${ENDPOINT}/${mail.id}`, mail),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mail"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    },
  })
}

export function useMailDelete() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id }: Pick<Mail, "id">) => apiDelete(`${ENDPOINT}/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mail"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    }
  })
}

export function useMailResend() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id }: Pick<Mail, "id">) => apiPost(`${ENDPOINT}/resend/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mail"] });
      queryClient.invalidateQueries({ queryKey: ["event"] });
    }
  })
}
