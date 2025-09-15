import { useQuery } from "@tanstack/react-query";
import { convertYearsToModel } from "../types/year";
import { apiGet } from "./query";

const ENDPOINT = "year";
const STALE_MIN_60 = 60 * 60 * 1000;

export function useYearGetAll(enabled: boolean = true) {
  return useQuery({
    queryKey: ["year"],
    queryFn: async () => (await apiGet(ENDPOINT, convertYearsToModel)).data,
    staleTime: STALE_MIN_60,
    throwOnError: true,
    enabled,
  });
}
