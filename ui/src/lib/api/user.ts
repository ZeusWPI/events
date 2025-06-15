import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertOrganizerToModel } from "../types/organizer";
import { apiGet, apiPost } from "./query";

const ENDPOINT_AUTH = "auth";
const ENDPOINT_USER = "organizer";

const STALE_30_MIN = 30 * 60 * 1000;

export function useUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertOrganizerToModel)).data,
    retry: 0,
    staleTime: STALE_30_MIN,
  });
}

// eslint-disable-next-line react-hooks-extra/no-unnecessary-use-prefix
export function useUserLogin() {
  window.location.href = `/api/${ENDPOINT_AUTH}/login/zauth`;
}

export function useUserLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => (await apiPost(`auth/logout`)).data,
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["user"] }),
  });
}
