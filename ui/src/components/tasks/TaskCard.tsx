import { useTaskStart } from "@/lib/api/task";
import type { Task } from "@/lib/types/task";
import { TaskStatus, TaskType } from "@/lib/types/task";
import { cn, formatDate } from "@/lib/utils/utils";
import { useNavigate } from "@tanstack/react-router";
import { Calendar, CalendarDays, LoaderCircle, Play } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Pill } from "../atoms/Pill";
import { TooltipText } from "../atoms/TooltipText";
import { Button } from "../ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";

interface Props {
  task: Task;
}

export function TaskCard({ task }: Props) {
  const navigate = useNavigate();
  const startTask = useTaskStart();

  const [updating, setUpdating] = useState(false);

  const recurring = task.type === TaskType.Recurring
  const running = task.status === TaskStatus.Running

  const handleCard = () => {
    if (!recurring) {
      return
    }

    void navigate({ to: "/tasks/$uid", params: { uid: task.uid } });
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
    <div onClick={handleCard} className={cn("flex justify-between items-center w-full p-4 bg-accent rounded-lg border ", recurring && "transition-transform duration-300 hover:scale-102 hover:cursor-pointer")}>
      <div className="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger>
            {recurring ? <CalendarDays /> : <Calendar />}
          </TooltipTrigger>
          <TooltipContent>
            {recurring ? "Recurring Task" : "One time Task"}
          </TooltipContent>
        </Tooltip>
        <div className="flex flex-col">
          <div className="flex align-center gap-2">
            <span>{task.name}</span>
            {task.lastError && (
              <Pill color="red">
                Failed
              </Pill>
            )}
            {running && (
              <Pill color="green">
                Running
              </Pill>
            )}
          </div>
          <span className="text-muted-foreground text-sm">{`Next run: ${formatDate(task.nextRun)}`}</span>
        </div>
      </div>
      <div className="flex items-center gap-2">
        {recurring && (
          <TooltipText text={running ? "Running" : "Run"}>
            <Button onClick={handleRun} disabled={running || updating} size="icon" className="rounded-full">
              {running ? <LoaderCircle className="animate-spin" /> : <Play />}
            </Button>
          </TooltipText>
        )}
      </div>
    </div>
  );
}
