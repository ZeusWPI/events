import type { TaskHistoryFilter } from "@/lib/types/task";
import { Calendar, CalendarDays, CircleX } from "lucide-react";
import { useState } from "react";
import useInfiniteScroll from "react-infinite-scroll-hook";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { BottomOfPage } from "@/components/molecules/BottomOfPage";
import { PageHeader } from "@/components/molecules/PageHeader";
import { TaskHistoryTable } from "@/components/tasks/TaskHistoryTable";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { useTaskGetHistory } from "@/lib/api/task";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";

export function TasksHistory() {
  const [filters, setFilters] = useState<TaskHistoryFilter>({ onlyErrored: false, recurring: undefined });
  const { history, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useTaskGetHistory(filters);

  useBreadcrumb({ title: "History", link: { to: "/tasks/history" } });

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  if (isLoading) {
    return <Indeterminate />;
  }

  return (
    <div className="space-y-8">
      <PageHeader>
        <Title>
          History
        </Title>
      </PageHeader>
      <div className="grid grid-cols-3 gap-4">
        <Card
          onClick={() => setFilters(val => ({ ...val, onlyErrored: !val.onlyErrored }))}
          className="cursor-pointer"
        >
          <CardContent>
            <div className="flex flex-col gap-2">
              <div className="flex justify-between">
                <CircleX />
                <Checkbox checked={filters.onlyErrored} />
              </div>
              <span>Errored</span>
            </div>
          </CardContent>
        </Card>
        <Card
          onClick={() => setFilters(val => ({ ...val, recurring: val.recurring !== false ? false : undefined }))}
          className="cursor-pointer"
        >
          <CardContent>
            <div className="flex flex-col gap-2">
              <div className="flex justify-between">
                <CalendarDays />
                <Checkbox checked={filters.recurring !== false} />
              </div>
              <span>Recurring</span>
            </div>
          </CardContent>
        </Card>
        <Card
          onClick={() => setFilters(val => ({ ...val, recurring: val.recurring !== true ? true : undefined }))}
          className="cursor-pointer"
        >
          <CardContent>
            <div className="flex flex-col gap-2">
              <div className="flex justify-between">
                <Calendar />
                <Checkbox checked={filters.recurring !== true} />
              </div>
              <span>One Time</span>
            </div>
          </CardContent>
        </Card>
      </div>
      <TaskHistoryTable history={history} />
      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </div>
  );
}
