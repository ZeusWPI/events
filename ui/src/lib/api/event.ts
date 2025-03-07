import type { Event, Year } from "../types/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertEventsToModel, convertEventToJSON } from "../utils/converter";
import { getApi, postApi } from "../utils/query";

const ENDPOINT = "event";

export function useEventByYear({ id }: Year) {
  return useQuery({
    queryKey: ["event", id],
    queryFn: async () => getApi(`${ENDPOINT}/year/${id}`, convertEventsToModel),
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
