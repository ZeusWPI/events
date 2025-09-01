import { useTaskResolve } from "@/lib/api/task";
import { TaskResult, type TaskHistory } from "@/lib/types/task";
import { formatDate } from "@/lib/utils/utils";
import type { ColumnDef, Row } from "@tanstack/react-table";
import { CalendarDaysIcon, CalendarIcon, ChevronDownIcon, ChevronUpIcon, FireExtinguisherIcon } from "lucide-react";
import { useMemo } from "react";
import { toast } from "sonner";
import { TooltipText } from "../atoms/TooltipText";
import { VirtualTable } from "../organisms/VirtualTable";
import { Button } from "../ui/button";

interface Props {
  history?: TaskHistory[];
  emptyText?: string;
}

const defaultHistory: TaskHistory[] = [];

function FormatDuration(nanos: number) {
  const ms = Math.floor(nanos / 1_000_000) % 1000
  const s = Math.floor(nanos / 1_000_000_000)

  const msString = `${ms.toString().padStart(3, '0')}ms`
  const sString = `${s.toString().padStart(2, '0')}s `

  return <span>{sString}<span className="text-muted-foreground">{msString}</span></span>
}

export function TaskHistoryTable({ history = defaultHistory, emptyText = "No history data yet" }: Props) {
  const columns: ColumnDef<TaskHistory>[] = useMemo(() => [
    {
      id: "status",
      cell: TaskIcon,
      meta: { small: true, horizontalAlign: "center" },
    },
    {
      accessorKey: "name",
      header: () => <span>Name</span>,
      cell: ({ cell }) => <span>{cell.getValue<string>()}</span>,
    },
    {
      accessorKey: "runAt",
      header: () => <span>Run at</span>,
      cell: row => <span>{formatDate(row.getValue<Date>())}</span>,
    },
    {
      accessorKey: "duration",
      header: () => <span>Duration</span>,
      cell: row => <span>{FormatDuration(row.getValue<number>())}</span>,
    },
    {
      id: "actions",
      cell: TaskActions,
      meta: { small: true, horizontalAlign: "center" },
    },
  ], []);

  if (!history.length) {
    return (
      <div className="text-center w-full">
        <span>{emptyText}</span>
      </div>
    );
  }

  return (
    <VirtualTable
      data={history}
      columns={columns}
      hasError={(item: TaskHistory) => item.error !== undefined}
      getError={(item: TaskHistory) => item.error ?? ""}
    />
  );
}

function TaskIcon({ row }: { row: Row<TaskHistory> }) {
  const task = row.original

  const color = () => {
    switch (task.result) {
      case TaskResult.FAILED:
        return "text-red-500"
      case TaskResult.RESOLVED:
        return "stroke-primary"
      default:
        return ""
    }
  }

  const text = () => {
    switch (task.result) {
      case TaskResult.FAILED:
        return "Failed"
      case TaskResult.RESOLVED:
        return "Resolved"
      default:
        return "Success"
    }
  }

  return (
    <TooltipText text={text()}>
      {task.recurring
        ? <CalendarDaysIcon className={color()} />
        : <CalendarIcon className={color()} />
      }
    </TooltipText>
  )
}

function TaskActions({ row }: { row: Row<TaskHistory> }) {
  const task = row.original
  const taskResolve = useTaskResolve()

  if (!row.getCanExpand()) {
    return null;
  }

  const handleResolve = () => {
    taskResolve.mutate(task, {
      onSuccess: () => toast.success(`Resolved ${task.name}`),
      onError: (error: Error) => toast.error("Failed", { description: error.message }),
    })
  }

  return (
    <div className="flex justify-center space-x-0">
      {task.result === TaskResult.FAILED && (
        <TooltipText text={"Mark as resolved"}>
          <Button onClick={handleResolve} size="icon" variant="ghost">
            <FireExtinguisherIcon />
          </Button>
        </TooltipText>
      )}
      <TooltipText text={row.getIsExpanded() ? "Hide error" : "Show error"}>
        <Button onClick={row.getToggleExpandedHandler()} size="icon" variant="ghost">
          {row.getIsExpanded() ? <ChevronUpIcon /> : <ChevronDownIcon />}
        </Button>
      </TooltipText>
    </div>
  )
}
