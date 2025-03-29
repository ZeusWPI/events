import type { TaskHistory } from "@/lib/types/task";
import type { ColumnDef, ExpandedState } from "@tanstack/react-table";
import { cn, formatDate } from "@/lib/utils/utils";
import { flexRender, getCoreRowModel, getExpandedRowModel, useReactTable } from "@tanstack/react-table";
import { Calendar, CalendarDays, ChevronDown, ChevronUp } from "lucide-react";
import React, { useMemo, useState } from "react";
import { Button } from "../ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";

interface Props {
  history?: TaskHistory[];
  emptyText?: string;
}

const defaultHistory: TaskHistory[] = [];

export function TaskHistoryTable({ history = defaultHistory, emptyText = "No history data yet" }: Props) {
  const [expanded, setExpanded] = useState<ExpandedState>({});

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
          <div className="flex justify-end w-full">
            <Button
              onClick={row.getToggleExpandedHandler()}
              size="icon"
              variant="ghost"
            >
              <Tooltip>
                <TooltipTrigger>
                  {row.getIsExpanded() ? <ChevronUp /> : <ChevronDown />}
                </TooltipTrigger>
                <TooltipContent>
                  {row.getIsExpanded() ? "Hide error" : "Show error"}
                </TooltipContent>
              </Tooltip>
            </Button>
          </div>
        );
      },
    },
  ], []);

  const table = useReactTable({
    data: history,
    columns,
    state: {
      expanded,
    },
    getCoreRowModel: getCoreRowModel(),
    getRowCanExpand: row => row.original.error !== undefined,
    getExpandedRowModel: getExpandedRowModel(),
    onExpandedChange: setExpanded,
  });

  const { getHeaderGroups, getRowModel } = table;
  const headerGroups = getHeaderGroups();
  const { rows } = getRowModel();

  if (!history.length) {
    return (
      <div className="text-center w-full">
        <span>{emptyText}</span>
      </div>
    );
  }

  return (
    <table className="relative w-full border-separate border-spacing-0">
      <thead>
        {headerGroups.map(headerGroup => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map(
              (header, idx) =>
                !header.column.columnDef.meta?.hideHeader && (
                  <th
                    key={header.id}
                    colSpan={header.column.columnDef.meta?.colSpanHeader ?? header.colSpan}
                    className={cn(
                      "p-3 text-sm font-bold bg-muted text-start sticky top-1",
                      idx === 0 && "rounded-l-md",
                      idx === headerGroup.headers.length - 1 && "rounded-r-md",
                      header.column.columnDef.meta?.small && "w-0",
                    )}
                  >
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </th>
                ),
            )}
          </tr>
        ))}
      </thead>
      <tbody>
        {rows.map(row => (
          <React.Fragment key={row.id}>
            <tr>
              {row.getVisibleCells().map(cell => (
                <td
                  key={cell.id}
                  className={cn(
                    "p-2",
                    !row.getIsExpanded() && "border-b",
                  )}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
            {row.getIsExpanded() && (
              <tr>
                <td colSpan={row.getAllCells().length} className="p-2 border-b">
                  <span className="text-red-400">{row.original.error}</span>
                </td>
              </tr>
            )}
          </React.Fragment>
        ))}
      </tbody>
    </table>
  );
}
