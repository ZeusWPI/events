import type { ColumnDef } from "@tanstack/react-table";
import type { TaskHistory } from "@/lib/types/task";
import { Calendar, CalendarDays, ChevronDown, ChevronUp } from "lucide-react";
import { cn, formatDate } from "@/lib/utils/utils";
import { Button } from "../ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";
import { Table } from "../organisms/Table";
import { useMemo } from "react";

interface Props {
  history?: TaskHistory[];
  emptyText?: string;
}

const defaultHistory: TaskHistory[] = [];

export function TaskHistoryTable({ history = defaultHistory, emptyText = "No history data yet" }: Props) {
  const columns: ColumnDef<TaskHistory>[] = useMemo(() => [
    {
      id: "recurring",
      cell: ({ row }) => (
        row.original.recurring
          ? <CalendarDays className={cn(row.original.error && "text-red-500")} />
          : <Calendar className={cn(row.original.error && "text-red-500")} />
      ),
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
      id: "actions",
      cell: ({ row }) => {
        if (!row.getCanExpand()) {
          return null;
        }

        return (
          <Button
            onClick={row.getToggleExpandedHandler()}
            size="icon"
            variant="ghost"
          >
            <Tooltip>
              <TooltipTrigger asChild>
                {row.getIsExpanded() ? <ChevronUp /> : <ChevronDown />}
              </TooltipTrigger>
              <TooltipContent>
                {row.getIsExpanded() ? "Hide error" : "Show error"}
              </TooltipContent>
            </Tooltip>
          </Button>
        );
      },
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
    <Table
      data={history}
      columns={columns}
      hasError={(item: TaskHistory) => item.error !== undefined}
      getError={(item: TaskHistory) => item.error ?? ""}
    />
  );
}
