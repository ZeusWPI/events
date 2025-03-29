import type { Event } from "../types/event";
import type { Year } from "../types/year";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertEventsToModel, convertEventToJSON } from "../types/event";
import { getApi, postApi } from "../utils/query";

const ENDPOINT = "event";
const MIN_5 = 5 * 60 * 1000;

export function useEventByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["event", id],
    queryFn: async () => getApi(`${ENDPOINT}/year/${id}`, convertEventsToModel),
    staleTime: MIN_5,
  });
}

export function useEventSaveOrganizers() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (events: Event[]) => postApi(`${ENDPOINT}/organizers`, events.map(event => convertEventToJSON(event))),
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
