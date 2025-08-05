import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { apiDelete, apiGet, apiPost, apiPut, NO_CONVERTER } from "./query"
import { CONTENT_TYPE } from "../types/contentType"
import { getUuid } from "../utils/utils"
import { Poster } from "../types/poster"
import { Year } from "../types/year"

const ENDPOINT = "poster"
const STALE_30_MIN = 30 * 60 * 1000

export function usePosterGetFile(poster: Pick<Poster, 'id' | 'eventId'>) {
  return useQuery({
    queryKey: ["event", poster.eventId, "poster", poster.id],
    queryFn: async () => {
      const { data } = await apiGet<Blob>(`${ENDPOINT}/${poster.id}/file`)
      return new File([data], `${getUuid()}.png`, { type: CONTENT_TYPE.PNG })
    },
    staleTime: STALE_30_MIN,
    enabled: poster.id > 0,
    throwOnError: true,
  })
}

export function usePosterCreate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (args: { poster: Poster, file: File }) => apiPut(ENDPOINT, args.poster, NO_CONVERTER, [{ file: args.file, field: "file" }]),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  })
}

export function usePosterUpdate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (args: { poster: Poster, file: File }) => apiPost(`${ENDPOINT}/${args.poster.id}`, args.poster, NO_CONVERTER, [{ file: args.file, field: "file" }]),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] }),
  })
}

export function usePosterDelete() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (args: { poster: Poster, year: Pick<Year, 'id'> }) => apiDelete(`${ENDPOINT}/${args.poster.id}`),
    onSuccess: (_, args) => queryClient.invalidateQueries({ queryKey: ["event", args.year.id] })
  })
}


