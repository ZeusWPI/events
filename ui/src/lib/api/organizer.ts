import type { Year } from "../types/types";
import { useQuery } from "@tanstack/react-query";
import { convertOrganizersToModel } from "../utils/converter";
import { getApi } from "../utils/query";

const ENDPOINT = "organizer";

export function useOrganizerByYear({ id }: Year) {
  return useQuery({
    queryKey: ["organizer", id],
    queryFn: async () => getApi(`${ENDPOINT}/year/${id}`, convertOrganizersToModel),
  });
}
