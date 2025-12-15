import type { Task, TaskHistory, TaskHistoryFilter } from "../types/task";
import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertTaskHistoryToModel, convertTasksToModel } from "../types/task";
import { apiGet, apiPost } from "./query";

const ENDPOINT = "task";
const PAGE_LIMIT = 100;
const STALE_MIN_5 = 5 * 60 * 1000;
const REFETCH_SEC_10 = 10 * 1000;

export function useTaskGetAll() {
  const queryClient = useQueryClient();

  return useQuery({
    queryKey: ["task"],
    queryFn: async () => (await apiGet(ENDPOINT, convertTasksToModel)).data,
    refetchInterval: REFETCH_SEC_10,
    structuralSharing(oldData, newData) {
      if (JSON.stringify(oldData) !== JSON.stringify(newData)) {
        void queryClient.invalidateQueries({ queryKey: ["task_history"] });
      }

      return newData;
    },
    throwOnError: true,
  });
}

export function useTaskGetHistory(filter?: TaskHistoryFilter) {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["task_history", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      });

      if (filter?.uid !== undefined) {
        queryParams.append("uid", filter.uid);
      }

      if (filter?.result !== undefined) {
        queryParams.append("result", filter.result.toString())
      }

      const url = `${ENDPOINT}/history?${queryParams.toString()}`;
      return (await apiGet(url, convertTaskHistoryToModel)).data;
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1;
    },
    enabled: filter !== undefined,
    staleTime: STALE_MIN_5,
    throwOnError: true,
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

export function useTaskResolve() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id }: Pick<TaskHistory, "id">) => apiPost(`${ENDPOINT}/resolve/${id}`),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["task"] });
      void queryClient.invalidateQueries({ queryKey: ["task_history"] });
    },
  });
}

export function useTaskStart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ uid }: Pick<Task, "uid">) => apiPost(`${ENDPOINT}/start/${uid}`),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["task"] });
      void queryClient.invalidateQueries({ queryKey: ["task_history"] });
    },
  });
}
