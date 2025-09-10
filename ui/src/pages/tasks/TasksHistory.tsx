import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { BottomOfPage } from "@/components/molecules/BottomOfPage";
import { PageHeader } from "@/components/molecules/PageHeader";
import { TaskHistoryTable } from "@/components/tasks/TaskHistoryTable";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { useTaskGetHistory } from "@/lib/api/task";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { weightSubcategory } from "@/lib/types/general";
import { TaskResult, type TaskHistoryFilter } from "@/lib/types/task";
import { CheckIcon, FireExtinguisherIcon, FlameIcon } from "lucide-react";
import { useState } from "react";
import useInfiniteScroll from "react-infinite-scroll-hook";

export function TasksHistory() {
  const [filter, setFilter] = useState<TaskHistoryFilter>({});
  const { history, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useTaskGetHistory(filter);

  useBreadcrumb({ title: "History", weight: weightSubcategory, link: { to: "/tasks/history" } });

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  if (isLoading) {
    return <Indeterminate />;
  }

  const handleClick = (newResult: TaskResult) => {
    setFilter(val => ({ ...val, result: val.result === newResult ? undefined : newResult }))
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
          onClick={() => handleClick(TaskResult.Failed)}
          className="cursor-pointer"
        >
          <CardContent>
            <div className="flex flex-col gap-2">
              <div className="flex justify-between">
                <FlameIcon />
                <Checkbox checked={filter.result === TaskResult.Failed} />
              </div>
              <span>Failed</span>
            </div>
          </CardContent>
        </Card>
        <Card
          onClick={() => handleClick(TaskResult.Resolved)}
          className="cursor-pointer"
        >
          <CardContent>
            <div className="flex flex-col gap-2">
              <div className="flex justify-between">
                <FireExtinguisherIcon />
                <Checkbox checked={filter.result === TaskResult.Resolved} />
              </div>
              <span>Resolved</span>
            </div>
          </CardContent>
        </Card>
        <Card
          onClick={() => handleClick(TaskResult.Succes)}
          className="cursor-pointer"
        >
          <CardContent>
            <div className="flex flex-col gap-2">
              <div className="flex justify-between">
                <CheckIcon />
                <Checkbox checked={filter.result === TaskResult.Succes} />
              </div>
              <span>Success</span>
            </div>
          </CardContent>
        </Card>
      </div>
      <TaskHistoryTable history={history} />
      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </div>
  );
}
