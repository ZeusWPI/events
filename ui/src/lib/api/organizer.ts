import type { Year } from "../types/year";
import { useQuery } from "@tanstack/react-query";
import { convertOrganizersToModel } from "../types/organizer";
import { apiGet } from "./query";

const ENDPOINT = "organizer";
const MIN_5 = 5 * 60 * 1000;

export function useOrganizerByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["organizer", id],
    queryFn: async () => (await apiGet(`${ENDPOINT}/year/${id}`, convertOrganizersToModel)).data,
    staleTime: MIN_5,
    throwOnError: true,
    enabled: id > 0,
  });
}
