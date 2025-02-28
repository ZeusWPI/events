import type { Year } from "../types/types";
import { useQuery } from "@tanstack/react-query";
import { convertEvents } from "../utils/converter";
import { getApi } from "../utils/query";

const ENDPOINT = "event";

export function useEventByYear({ id }: Year) {
  return useQuery({
    queryKey: ["event", id],
    queryFn: async () => getApi(`${ENDPOINT}/year/${id}`, convertEvents),
  });
}
