import type { Year } from "../types/year";
import { useQuery } from "@tanstack/react-query";
import { convertOrganizersToModel } from "../types/organizer";
import { getApi } from "../utils/query";

const ENDPOINT = "organizer";
const MIN_5 = 5 * 60 * 1000;

export function useOrganizerByYear({ id }: Pick<Year, "id">) {
  return useQuery({
    queryKey: ["organizer", id],
    queryFn: async () => getApi(`${ENDPOINT}/year/${id}`, convertOrganizersToModel),
    staleTime: MIN_5,
  });
}
