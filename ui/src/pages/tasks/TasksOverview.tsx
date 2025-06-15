import { Link } from "@tanstack/react-router";
import useInfiniteScroll from "react-infinite-scroll-hook";
import { DividerText } from "@/components/atoms/DividerText";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { BottomOfPage } from "@/components/molecules/BottomOfPage";
import { PageHeader } from "@/components/molecules/PageHeader";
import { TaskCard } from "@/components/tasks/TaskCard";
import { TaskHistoryTable } from "@/components/tasks/TaskHistoryTable";
import { Button } from "@/components/ui/button";
import { useTaskGetAll, useTaskGetHistory } from "@/lib/api/task";

export function TasksOverview() {
  const { data: tasks = [], isLoading } = useTaskGetAll();
  const { history, isFetchingNextPage, hasNextPage, fetchNextPage } = useTaskGetHistory({ onlyErrored: true });

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
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>Tasks</Title>
        <Button variant="outline" asChild>
          <Link to="/tasks/history">
            History
          </Link>
        </Button>
      </PageHeader>
      <div className="grid lg:grid-cols-2 gap-4">
        {tasks.map(task => (
          <div key={task.name} className="">
            <TaskCard task={task} />
          </div>
        ))}
      </div>
      <DividerText>
        Failed tasks
      </DividerText>
      <TaskHistoryTable history={history} emptyText="No failed tasks yet" />
      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </div>
  );
}
