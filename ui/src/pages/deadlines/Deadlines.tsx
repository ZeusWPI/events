import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Countdown } from "@/components/molecules/Countdown";
import { Table } from "@/components/organisms/Table";
import { useCheckGetByYear } from "@/lib/api/check";
import { useEventByYear } from "@/lib/api/event";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { useYear } from "@/lib/hooks/useYear";
import { Check, CheckStatus } from "@/lib/types/check";
import { Event } from "@/lib/types/event";
import { weightCategory } from "@/lib/types/general";
import { ColumnDef, Row } from "@tanstack/react-table";
import { useMemo } from "react";

interface Deadline {
  check: Check;
  event: Event;
}

export function Deadlines() {
  const { year } = useYear()
  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year)
  const { data: checks, isLoading: isLoadingChecks } = useCheckGetByYear(year)

  useBreadcrumb({ title: "Deadlines", weight: weightCategory, link: { to: "/deadlines" } })

  const columns: ColumnDef<Deadline>[] = useMemo(() => [
    {
      id: "name",
      header: () => <span>Check</span>,
      cell: ({ row }) => <span>{row.original.check.description}</span>
    },
    {
      id: "event",
      header: () => <span>Event</span>,
      cell: ({ row }) => <span>{row.original.event.name}</span>
    },
    {
      id: "deadline",
      header: () => <span>Deadline</span>,
      cell: ({ row }) => <Deadline row={row} />
    }
  ], [])

  const sorted = useMemo(() => sortChecks(checks ?? [], events ?? []), [checks, events])

  if (isLoadingEvents || isLoadingChecks) {
    return <Indeterminate />
  }

  return (
    <Table
      data={sorted}
      columns={columns}
    />
  )
}

function Deadline({ row }: { row: Row<Deadline> }) {
  const event = row.original.event
  const check = row.original.check

  const now = new Date()
  const deadline = new Date(event.startTime.getTime() - check.deadline!)
  const tooLate = deadline < now

  if (tooLate) {
    return null
  }

  return <Countdown goalDate={deadline} />
}

function sortChecks(allChecks: Check[], events: Event[]): Deadline[] {
  const now = new Date()
  const eventMap = new Map(events?.map(e => [e.id, e]) ?? {})
  const checks = allChecks?.filter(c => {
    if (c.deadline === undefined) return false
    if (!eventMap.has(c.eventId)) return false
    if (![CheckStatus.Todo, CheckStatus.Warning].includes(c.status)) return false

    const deadline = new Date(eventMap.get(c.eventId)!.startTime.getTime() - c.deadline)
    return deadline > now
  }) ?? []

  return [...checks].sort((a, b) => {
    const aDeadline = eventMap.get(a.eventId)!.startTime.getTime() - a.deadline!
    const bDeadline = eventMap.get(b.eventId)!.startTime.getTime() - b.deadline!

    return aDeadline - bDeadline
  }).map(c => ({ check: c, event: eventMap.get(c.eventId)! }))
}
