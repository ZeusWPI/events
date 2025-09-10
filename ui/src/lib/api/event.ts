import type { Event } from "../types/event";
import type { Year } from "../types/year";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertEventsToModel, convertEventToModel } from "../types/event";
import { apiGet, apiPost } from "./query";

const ENDPOINT = "event";
const STALE_5_MIN = 5 * 60 * 1000;

export function useEvent({ id }: Pick<Event, "id">) {
  return useQuery({
    queryKey: ["event", id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/${id}`, convertEventToModel)).data,
    staleTime: STALE_5_MIN,
    throwOnError: true,
    enabled: id > 0,
  })
}

export function useEventByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["event", "year", id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/year/${id}`, convertEventsToModel)).data,
    staleTime: STALE_5_MIN,
    throwOnError: true,
    enabled: id > 0,
  });
}

export function useEventSaveOrganizers() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (eventOrganizers: Pick<Event, 'id' | 'organizers' | 'year'>[]) => apiPost(`${ENDPOINT}/organizers`, eventOrganizers.map(e => ({ eventId: e.id, organizers: e.organizers.map(o => o.id) }))),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["event"] })
  });
}
