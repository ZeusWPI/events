import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertOrganizerToModel } from "../utils/converter";
import { getApi, postApi } from "../utils/query";

export function useUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => getApi(`organizer/me`, convertOrganizerToModel),
    retry: 0,
  });
}

export function useUserLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => postApi(`auth/logout`),
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["user"] }),
  });
}
