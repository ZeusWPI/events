import type { Task, TaskHistoryFilter } from "../types/task";
import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertTaskHistoryToModel, convertTasksToModel } from "../types/task";
import { apiGet, apiPost } from "./query";

const ENDPOINT = "task";
const PAGE_LIMIT = 100;
const STALE_MIN_5 = 5 * 60 * 1000;
const REFETCH_SEC_30 = 10 * 1000;

export function useTaskGetAll() {
  const queryClient = useQueryClient();

  return useQuery({
    queryKey: ["task"],
    queryFn: async () => (await apiGet(ENDPOINT, convertTasksToModel)).data,
    refetchInterval: REFETCH_SEC_30,
    structuralSharing(oldData, newData) {
      if (JSON.stringify(oldData) !== JSON.stringify(newData)) {
        void queryClient.invalidateQueries({ queryKey: ["task_history"] });
      }

      return newData;
    },
  });
}

export function useTaskGetHistory(filters?: TaskHistoryFilter) {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["task_history", filters],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      });

      if (filters?.name !== undefined) {
        queryParams.append("name", filters.name);
      }

      if (filters?.onlyErrored !== undefined) {
        queryParams.append("only_errored", filters.onlyErrored.toString());
      }

      if (filters?.recurring !== undefined) {
        queryParams.append("recurring", filters.recurring.toString());
      }

      const url = `${ENDPOINT}/history?${queryParams.toString()}`;
      return (await apiGet(url, convertTaskHistoryToModel)).data;
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1;
    },
    enabled: filters !== undefined,
    staleTime: STALE_MIN_5,
  });

  const history = data?.pages.flat() ?? [];

  return {
    history,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  };
}

export function useTaskStart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id }: Pick<Task, "id">) => apiPost(`${ENDPOINT}/${id}`),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["task"] });
      void queryClient.invalidateQueries({ queryKey: ["task_history"] });
    },
  });
}
