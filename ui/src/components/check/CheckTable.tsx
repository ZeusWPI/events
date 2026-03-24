import { useCheckCreate, useCheckDelete, useCheckDone, useCheckUpdate } from "@/lib/api/check";
import { Check, CheckStatus, CheckType, checkStatusToIcon, checkStatusToText } from "@/lib/types/check";
import { Event } from "@/lib/types/event";
import { ColumnDef } from "@tanstack/react-table";
import { ClipboardCheckIcon, ClipboardXIcon, MessageCircleMoreIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { IconButton } from "../atoms/IconButton";
import { TooltipText } from "../atoms/TooltipText";
import { Countdown } from "../molecules/Countdown";
import { Table } from "../organisms/Table";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Confirm } from "../molecules/Confirm";

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
      id: "deadline",
      header: () => null,
      cell: ({ row }) => row.original.type === CheckType.Automatic
        ? <AutomaticDeadline check={row.original} event={event} />
        : <ManualDeadline />
    },
    {
      id: "actions",
      header: () => (
        <div className={`flex justify-end ${addCheck ? "invisible" : "flex"}`}>
          <IconButton onClick={() => setAddCheck(true)}>
            <PlusIcon className="text-primary" />
          </IconButton>
        </div>
      ),
      cell: ({ row }) => row.original.type === CheckType.Automatic
        ? <AutomaticActions check={row.original} />
        : <ManualActions check={row.original} />
      ,
      meta: { small: true, horizontalAlign: "justify-end" },
    }
  ], [checks, addCheck]) // eslint-disable-line react-hooks/exhaustive-deps

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

interface ManualActionsProps {
  check: Check;
}

function ManualActions({ check }: ManualActionsProps) {
  const [toggleStatus, setToggleStatus] = useState(false)
  const CheckUpdate = useCheckUpdate()

  const toggleDone = () => {
    setToggleStatus(true)

    const newCheck = { ...check, status: check.status === CheckStatus.Todo ? CheckStatus.Done : CheckStatus.Todo }

    CheckUpdate.mutate(newCheck, {
      onSuccess: () => toast.success("Success"),
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setToggleStatus(false)
    })
  }

  const [deleteStatus, setDeleteStatus] = useState(false)
  const checkDelete = useCheckDelete()

  const deleteCheck = () => {
    setDeleteStatus(true)

    checkDelete.mutate(check, {
      onSuccess: () => toast.success("Success"),
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setDeleteStatus(false),
    })
  }

  return (
    <div className="flex">
      <Button onClick={toggleDone} size="iconSmall" variant="ghost" disabled={toggleStatus}>
        {check.status === CheckStatus.Done ? <ClipboardXIcon /> : <ClipboardCheckIcon />}
      </Button>
      <Button onClick={deleteCheck} size="iconSmall" variant="ghost" disabled={deleteStatus}>
        <Trash2Icon className="text-red-500" />
      </Button>
    </div>
  )
}

function ManualDeadline() {
  return null
}

interface AutomaticDeadlineProps {
  event: Event;
  check: Check;
}

function AutomaticDeadline({ event, check }: AutomaticDeadlineProps) {
  if (!check.deadline) {
    // Some automatic tasks have no deadline
    return
  }

  if ([CheckStatus.Done, CheckStatus.DoneLate].includes(check.status)) {
    return
  }

  const now = new Date()
  const deadline = new Date(event.startTime.getTime() - check.deadline)
  const tooLate = deadline < now

  if (tooLate) {
    return
  }

  return <Countdown goalDate={deadline} />
}

interface AutomaticActionsProps {
  check: Check;
}

function AutomaticActions({ check }: AutomaticActionsProps) {
  const checkDone = useCheckDone()

  const [open, setOpen] = useState(false)

  if (![CheckStatus.Todo, CheckStatus.TodoLate].includes(check.status)) return

  const handleDone = () => {
    setOpen(true)
  }

  const handleDoneConfirm = () => {
    checkDone.mutate(check, {
      onSuccess: () => {
        toast.success("Success")
      },
      onError: error => toast.error("Failed", { description: error.message }),
    })
  }

  return (
    <>
      <TooltipText text="Mark as done">
        <Button onClick={handleDone} size="iconSmall" variant="ghost">
          <ClipboardCheckIcon />
        </Button>
      </TooltipText>
      <Confirm
        title="Mark as done confirmation"
        description="Are you sure you want to mark this task as done"
        confirmText="Mark as done"
        onConfirm={handleDoneConfirm}
        open={open}
        onOpenChange={setOpen}
      />
    </>
  )
}
