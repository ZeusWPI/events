import { useParams } from "@tanstack/react-router";
import { LoaderCircle, Play } from "lucide-react";
import { useMemo, useState } from "react";
import useInfiniteScroll from "react-infinite-scroll-hook";
import { toast } from "sonner";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Pill } from "@/components/atoms/Pill";
import { Title } from "@/components/atoms/Title";
import { BottomOfPage } from "@/components/molecules/BottomOfPage";
import { PageHeader } from "@/components/molecules/PageHeader";
import { TaskHistoryTable } from "@/components/tasks/TaskHistoryTable";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useTaskGetAll, useTaskGetHistory, useTaskStart } from "@/lib/api/task";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { TaskStatus } from "@/lib/types/task";
import { formatDate } from "@/lib/utils/utils";
import Error404 from "../404";

interface Interval {
  day: number;
  hour: number;
  minute: number;
}

function convertNanoToInterval(nanoseconds: number): Interval {
  const totalMinutes = nanoseconds / (1e9 * 60);
  const day = Math.floor(totalMinutes / (24 * 60));
  const hour = Math.floor((totalMinutes % (24 * 60)) / 60);
  const minute = Math.floor(totalMinutes % 60);

  return { day, hour, minute };
}

export function TasksDetail() {
  const { id: taskID } = useParams({ from: "/tasks/$id" });

  const { data: tasks = [], isLoading } = useTaskGetAll();
  const task = tasks.find(task => task.id.toString() === taskID);
  const { history, isLoading: isHistoryLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useTaskGetHistory(task ? { name: task.name } : undefined);

  const startTask = useTaskStart();

  useBreadcrumb({ title: task?.name ?? "", link: { to: "/tasks/$id", params: { id: task?.id.toString() ?? "" } } });

  const [updating, setUpdating] = useState(false);

  const interval = useMemo(() => convertNanoToInterval(task?.interval ?? 0), [task?.interval]);

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  if (isLoading) {
    return <Indeterminate />;
  }

  if (!task || !task.recurring) {
    return <Error404 />;
  }

  const handleRun = () => {
    setUpdating(true);

    startTask.mutate(task, {
      onSuccess: () => {
        toast.success(task.name, { description: "Started" });
      },
      onError: () => toast.error(task.name, { description: "Failed to start" }),
      onSettled: () => setUpdating(false),
    });
  };

  return (
    <div className="space-y-8">
      <PageHeader className="col-span-full items-center">
        <div className="flex flex-col gap-1.5">
          <div className="flex align-center space-x-2">
            <Title>{task.name}</Title>
            {task.status === TaskStatus.running && (
              <Pill color="green">
                Running
              </Pill>
            )}
          </div>
          <div className="flex flex-col md:flex-row gap-0 md:gap-1.5 align-center text-muted-foreground text-base/4 h-full">
            <span>{`Next run: ${formatDate(task.nextRun)}`}</span>
            {task.lastRun && (
              <>
                <Separator orientation="vertical" className="hidden sm:block h-4 w-1" />
                <span>{`Last run: ${formatDate(task.lastRun)}`}</span>
              </>
            )}
          </div>
        </div>
        <Button onClick={handleRun} disabled={task.status === TaskStatus.running || updating} size="icon" className="rounded-full">
          {task.status === TaskStatus.running ? <LoaderCircle className="animate-spin" /> : <Play />}
        </Button>
      </PageHeader>
      <div className="flex justify-between align-center p-6 w-full border rounded-md gap-4">
        <div className="flex flex-col align-start grow-1">
          <span className="text-muted-foreground">Day</span>
          <span className="text-3xl text-center">{interval.day}</span>
        </div>
        <Separator orientation="vertical" className="h-15" />
        <div className="flex flex-col align-start grow-1">
          <span className="text-muted-foreground">Hour</span>
          <span className="text-3xl text-center">{interval.hour}</span>
        </div>
        <Separator orientation="vertical" className="h-15" />
        <div className="flex flex-col align-start grow-1">
          <span className="text-muted-foreground">Minute</span>
          <span className="text-3xl text-center">{interval.minute}</span>
        </div>
      </div>
      <div>
        {!isHistoryLoading
          ? <TaskHistoryTable history={history} />
          : <Indeterminate />}
      </div>
      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </div>
  );
}
