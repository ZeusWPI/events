import type { Event } from "../types/event";
import type { Year } from "../types/year";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertEventsToModel, convertEventToJSON } from "../types/event";
import { apiGet, apiPost } from "./query";

const ENDPOINT = "event";
const MIN_5 = 5 * 60 * 1000;

export function useEventByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["event", id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/year/${id}`, convertEventsToModel)).data,
    staleTime: MIN_5,
    throwOnError: true,
  });
}

export function useEventSaveOrganizers() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (events: Event[]) => apiPost(`${ENDPOINT}/organizers`, events.map(convertEventToJSON)),
    onSuccess: (_, events) => {
      const uniqueYears = Array.from(
        new Set(events.map(event => event.year.id)),
      );

      uniqueYears.forEach((yearId) => {
        void queryClient.invalidateQueries({ queryKey: ["event", yearId] });
      });
    },
  });
}
