import { Check, CheckSource } from "@/lib/types/check"
import { Tooltip, TooltipTrigger, TooltipContent } from "@radix-ui/react-tooltip";
import { ColumnDef } from "@tanstack/react-table";
import { CheckIcon, ChevronDownIcon, ChevronUpIcon, ClipboardCheckIcon, ClipboardXIcon, PlusIcon, XIcon } from "lucide-react";
import { useMemo, useState } from "react";
import { Button } from "../ui/button";
import { useCheckCreate, useCheckToggle } from "@/lib/api/check";
import { toast } from "sonner";
import { Input } from "../ui/input";
import { Table } from "../organisms/Table";

interface Props {
  checks: Check[];
  eventId: number;
}

export function CheckTable({ checks, eventId }: Props) {
  const [adding, setAdding] = useState(false)
  const [addCheck, setAddCheck] = useState(false)
  const [addCheckDescription, setAddCheckDescription] = useState("")

  const checkCreate = useCheckCreate()

  const cancelCheck = () => {
    setAddCheck(false)
    setAddCheckDescription("")
  }
  const createCheck = () => {
    if (addCheckDescription === "") {
      return
    }

    setAdding(true)

    checkCreate.mutate({ eventId: eventId, description: addCheckDescription }, {
      onSuccess: () => {
        toast.success("Success")
        setAddCheck(false)
        setAddCheckDescription("")
      },
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setAdding(false)
    })
  }

  const [toggleStatus, setToggleStatus] = useState(false)

  const checkToggle = useCheckToggle()

  const toggleDone = (check: Check) => {
    setToggleStatus(true)

    checkToggle.mutate(check, {
      onSuccess: () => toast.success("Success"),
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setToggleStatus(false)
    })
  }

  const columns: ColumnDef<Check>[] = useMemo(() => [
    {
      accessorKey: "done",
      header: () => <span>{`${checks.filter(c => c.done).length}/${checks.length}`}</span>,
      cell: ({ cell }) => {
        const done = cell.getValue<boolean>()

        if (done) return <CheckIcon className="text-green-500" />
        else return <XIcon className="text-red-500" />
      },
      meta: { small: true, horizontalAlign: "center" },
    },
    {
      accessorKey: "description",
      header: () => <span>Description</span>,
      cell: ({ cell }) => <span>{cell.getValue<string>()}</span>
    },
    {
      id: "actions",
      header: () => {
        if (addCheck) return null

        return (
          <Button onClick={() => setAddCheck(true)} size="icon" variant="outline">
            <PlusIcon className="text-primary" />
          </Button>
        )
      },
      cell: ({ row }) => {
        if (!row.getCanExpand()) {
          return null
        }

        const check: Check = row.original

        if (check.source === CheckSource.Automatic) {
          if (check.error) {
            <Button onClick={row.getToggleExpandedHandler()} size="icon" variant="outline">
              <Tooltip>
                <TooltipTrigger asChild>
                  {row.getIsExpanded() ? <ChevronUpIcon /> : <ChevronDownIcon />}
                </TooltipTrigger>
                <TooltipContent>
                  <span>{row.getIsExpanded() ? "Hide error" : "Show error"}</span>
                </TooltipContent>
              </Tooltip>
            </Button>
          }

          return null
        }

        return (
          <Button onClick={() => toggleDone(check)} size="icon" variant="outline" disabled={toggleStatus}>
            <Tooltip>
              <TooltipTrigger asChild>
                {check.done ? <ClipboardXIcon /> : <ClipboardCheckIcon />}
              </TooltipTrigger>
              <TooltipContent>
                <span>{check.done ? "Mark as undone" : "Mark as done"}</span>
              </TooltipContent>
            </Tooltip>
          </Button>
        )
      },
      meta: { small: true, horizontalAlign: "justify-end" },
    }
  ], [toggleStatus]) // eslint-disable-line react-hooks/exhaustive-deps
  return (
    <div className="space-y-4">
      {addCheck && (
        <div className="flex w-full gap-2">
          <Input placeholder="Description" onChange={e => setAddCheckDescription(e.target.value)} />
          <Button onClick={cancelCheck} variant="outline">
            <span>Cancel</span>
          </Button>
          <Button onClick={createCheck} disabled={adding}>
            <span>Create</span>
          </Button>
        </div>
      )}
      <Table
        data={checks}
        columns={columns}
        hasError={(item: Check) => item.error !== undefined}
        getError={(item: Check) => item.error ?? ""}
      />
    </div>
  );
}
