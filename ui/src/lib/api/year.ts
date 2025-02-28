import { useQuery } from "@tanstack/react-query";
import { convertYears } from "../utils/converter";
import { getApi } from "../utils/query";

const ENDPOINT = "year";

export function useYearGetAll() {
  return useQuery({
    queryKey: ["year"],
    queryFn: async () => getApi(ENDPOINT, convertYears),
  });
}
