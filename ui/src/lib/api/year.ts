import { useQuery } from "@tanstack/react-query";
import { convertYearsToModel } from "../types/year";
import { getApi } from "../utils/query";

const ENDPOINT = "year";
const STALE_MIN_60 = 60 * 60 * 1000;

export function useYearGetAll() {
  return useQuery({
    queryKey: ["year"],
    queryFn: async () => getApi(ENDPOINT, convertYearsToModel),
    staleTime: STALE_MIN_60,
  });
}
