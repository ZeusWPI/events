import { cn } from "@/lib/utils/utils";
import { useReactTable, getCoreRowModel, getExpandedRowModel, flexRender, ColumnDef, ExpandedState } from "@tanstack/react-table";
import React, { useState } from "react";
import { useVirtualizer } from '@tanstack/react-virtual'

interface Props<T> {
  data: T[];
  columns: ColumnDef<T>[];
  hasError?: (item: T) => boolean;
  getError?: (item: T) => string;
}

export function Table<T>({ data, columns, hasError, getError }: Props<T>) {
  const parentRef = React.useRef<HTMLDivElement>(null)
  const [expanded, setExpanded] = useState<ExpandedState>({});

  const table = useReactTable({
    data,
    columns,
    state: {
      expanded,
    },
    getCoreRowModel: getCoreRowModel(),
    getRowCanExpand: row => hasError?.(row.original) ?? false,
    getExpandedRowModel: getExpandedRowModel(),
    onExpandedChange: setExpanded,
  });

  const { getHeaderGroups, getRowModel } = table;
  const headerGroups = getHeaderGroups();
  const { rows } = getRowModel();

  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 34,
    overscan: 20,
  })

  return (
    <div ref={parentRef}>
      <div style={{ height: `${virtualizer.getTotalSize()}px` }}>
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
                          "p-2 text-sm font-bold bg-muted text-start sticky top-1",
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
            {virtualizer.getVirtualItems().map(virtualRow => {
              const row = rows[virtualRow.index]!
              return (
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
                        <span className="text-red-400">{getError?.(row.original)}</span>
                      </td>
                    </tr>
                  )}
                </React.Fragment>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
