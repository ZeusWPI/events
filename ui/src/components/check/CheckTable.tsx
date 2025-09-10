import { useCheckCreate, useCheckDelete, useCheckToggle } from "@/lib/api/check";
import { Check, CheckStatus, CheckType, checkStatusToIcon } from "@/lib/types/check";
import { ColumnDef, Row } from "@tanstack/react-table";
import { ChevronDownIcon, ChevronUpIcon, ClipboardCheckIcon, ClipboardXIcon, MessageSquareIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { IconButton } from "../atoms/IconButton";
import { TooltipText } from "../atoms/TooltipText";
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
      header: () => <span>{`${checks.filter(c => c.status === CheckStatus.Done || c.status === CheckStatus.DoneLate).length}/${checks.length}`}</span>,
      cell: ({ row }) => {
        if (!row.original.message) return checkStatusToIcon[row.original.status]

        return (
          <TooltipText text={row.original.message}>
            {checkStatusToIcon[row.original.status]}
            <MessageSquareIcon />
          </TooltipText>
        )
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
      cell: ({ row }) => <CheckActions row={row} onToggle={toggleDone} toggleStatus={toggleStatus} onDelete={deleteCheck} deleteStatus={deleteStatus} />,
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
      />
    </div>
  );
}

interface ActionProps {
  row: Row<Check>;
  onToggle: (check: Check) => void;
  toggleStatus: boolean;
  onDelete: (check: Check) => void;
  deleteStatus: boolean;
}
function CheckActions({ row, onToggle, toggleStatus, onDelete, deleteStatus }: ActionProps) {
  if (!row.getCanExpand()) {
    return null
  }

  const check: Check = row.original

  if (check.source === CheckType.Automatic) {
    if (check.error) {
      <IconButton onClick={row.getToggleExpandedHandler()}>
        {row.getIsExpanded() ? <ChevronUpIcon /> : <ChevronDownIcon />}
      </IconButton>
    }

    return null
  }

  return (
    <div className="flex">
      <Button onClick={() => onToggle(check)} size="icon" variant="ghost" disabled={toggleStatus}>
        {check.status === CheckStatus.Done ? <ClipboardXIcon /> : <ClipboardCheckIcon />}
      </Button>
      <Button onClick={() => onDelete(check)} size="icon" variant="ghost" disabled={deleteStatus}>
        <Trash2Icon className="text-red-500" />
      </Button>
    </div>
  )
}
