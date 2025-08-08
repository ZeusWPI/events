import type { Task } from "@/lib/types/task";
import { Link, useNavigate } from "@tanstack/react-router";
import { Calendar, CalendarDays, LoaderCircle, Play } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { useTaskStart } from "@/lib/api/task";
import { TaskStatus } from "@/lib/types/task";
import { cn, formatDate } from "@/lib/utils/utils";
import { Pill } from "../atoms/Pill";
import { Button } from "../ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";

interface Props {
  task: Task;
}

export function TaskCard({ task }: Props) {
  const navigate = useNavigate();
  const startTask = useTaskStart();

  const [updating, setUpdating] = useState(false);

  const handleCard = () => {
    if (!task.recurring) {
      return
    }

    void navigate({ to: "/tasks/$id", params: { id: task.id.toString() } });
  }

  const handleRun = (e: React.MouseEvent) => {
    e.stopPropagation(); // Avoid triggering the card press

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
    <div onClick={handleCard} className={cn("flex justify-between items-center w-full p-4 bg-accent rounded-lg border ", task.recurring && "transition-transform duration-300 hover:scale-102 hover:cursor-pointer")}>
      <div className="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger>
            {task.recurring ? <CalendarDays /> : <Calendar />}
          </TooltipTrigger>
          <TooltipContent>
            {task.recurring ? "Recurring Task" : "One time Task"}
          </TooltipContent>
        </Tooltip>
        <div className="flex flex-col">
          <div className="flex align-center gap-2">
            <Link to="/tasks/$id" params={{ id: task.id.toString() }}>
              <span>{task.name}</span>
            </Link>
            {task.lastError && (
              <Pill color="red">
                Failed
              </Pill>
            )}
            {task.status === TaskStatus.running && (
              <Pill color="green">
                Running
              </Pill>
            )}
          </div>
          <span className="text-muted-foreground text-sm">{`Next run: ${formatDate(task.nextRun)}`}</span>
        </div>
      </div>
      <div className="flex items-center gap-2">
        {task.recurring && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button onClick={handleRun} disabled={task.status === TaskStatus.running || updating} size="icon" className="rounded-full">
                {task.status === TaskStatus.running ? <LoaderCircle className="animate-spin" /> : <Play />}
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <span>{task.status === TaskStatus.running ? "Running" : "Run"}</span>
            </TooltipContent>
          </Tooltip>
        )}
      </div>
    </div>
  );
}
