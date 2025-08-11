import { useCheckCreate, useCheckDelete, useCheckToggle } from "@/lib/api/check";
import { Check, CheckSource, CheckStatus, statusToIcon } from "@/lib/types/check";
import { ColumnDef } from "@tanstack/react-table";
import { ChevronDownIcon, ChevronUpIcon, ClipboardCheckIcon, ClipboardXIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { IconButton } from "../atoms/IconButton";
import { Table } from "../organisms/Table";
import { Button } from "../ui/button";
import { Input } from "../ui/input";

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

  const [deleteStatus, setDeleteStatus] = useState(false)
  const checkDelete = useCheckDelete()

  const deleteCheck = (check: Check) => {
    setDeleteStatus(true)

    checkDelete.mutate(check, {
      onSuccess: () => toast.success("Success"),
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setDeleteStatus(false),
    })
  }

  const columns: ColumnDef<Check>[] = useMemo(() => [
    {
      accessorKey: "status",
      header: () => <span>{`${checks.filter(c => c.status === CheckStatus.Finished).length}/${checks.length}`}</span>,
      cell: ({ cell }) => {
        const status = cell.getValue<CheckStatus>()
        return statusToIcon[status]
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
          <div className="flex justify-end">
            <IconButton onClick={() => setAddCheck(true)}>
              <PlusIcon className="text-primary" />
            </IconButton>
          </div>
        )
      },
      cell: ({ row }) => {
        if (!row.getCanExpand()) {
          return null
        }

        const check: Check = row.original

        if (check.source === CheckSource.Automatic) {
          if (check.error) {
            <IconButton onClick={row.getToggleExpandedHandler()}>
              {row.getIsExpanded() ? <ChevronUpIcon /> : <ChevronDownIcon />}
            </IconButton>
          }

          return null
        }

        return (
          <div className="flex">
            <Button onClick={() => toggleDone(check)} size="icon" variant="ghost" disabled={toggleStatus}>
              {check.status === CheckStatus.Finished ? <ClipboardXIcon /> : <ClipboardCheckIcon />}
            </Button>
            <Button onClick={() => deleteCheck(check)} size="icon" variant="ghost" disabled={deleteStatus}>
              <Trash2Icon className="text-red-500" />
            </Button>
          </div>
        )
      },
      meta: { small: true, horizontalAlign: "justify-end" },
    }
  ], [checks, toggleStatus]) // eslint-disable-line react-hooks/exhaustive-deps

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
