import { useCheckCreate, useCheckDelete, useCheckUpdate } from "@/lib/api/check";
import { Check, CheckStatus, CheckType, checkStatusToIcon, checkStatusToText } from "@/lib/types/check";
import { ColumnDef, Row } from "@tanstack/react-table";
import { ClipboardCheckIcon, ClipboardXIcon, MessageCircleMoreIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { IconButton } from "../atoms/IconButton";
import { TooltipText } from "../atoms/TooltipText";
import { Table } from "../organisms/Table";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Event } from "@/lib/types/event";
import { Countdown } from "../molecules/Countdown";

interface Props {
  checks: Check[];
  event: Event;
}

export function CheckTable({ checks, event }: Props) {
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
      toast.error("Empty description")
      return
    }

    setAdding(true)

    checkCreate.mutate({ eventId: event.id, description: addCheckDescription }, {
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
  const CheckUpdate = useCheckUpdate()

  const toggleDone = (check: Check) => {
    setToggleStatus(true)

    if (check.status === CheckStatus.Todo) {
      check.status = CheckStatus.Done
    } else {
      check.status = CheckStatus.Todo
    }

    CheckUpdate.mutate(check, {
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
      cell: ({ row }) => (
        <TooltipText text={checkStatusToText[row.original.status]} subtext={row.original.message}>
          <div className="flex gap-0">
            {checkStatusToIcon[row.original.status]}
            {row.original.message && <MessageCircleMoreIcon size={16} />}
          </div>
        </TooltipText>
      ),
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
      cell: ({ row }) => row.original.type === CheckType.Automatic
        ? <Deadline check={row.original} event={event} />
        : <CheckActions row={row} onToggle={toggleDone} toggleStatus={toggleStatus} onDelete={deleteCheck} deleteStatus={deleteStatus} />
      ,
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
  const check: Check = row.original

  if (check.type === CheckType.Automatic) {
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

interface DeadlineProps {
  event: Event;
  check: Check;
}

function Deadline({ event, check }: DeadlineProps) {
  if (check.type !== CheckType.Automatic) {
    return
  }

  if (!check.deadline) {
    // Some automatic tasks have no deadline
    return null
  }

  const now = new Date()
  const deadline = new Date(now.getTime() + (check.deadline / 1000))
  const tooLate = deadline > event.startTime

  if (tooLate) {
    return null
  }

  return <Countdown goalDate={deadline} />
}
