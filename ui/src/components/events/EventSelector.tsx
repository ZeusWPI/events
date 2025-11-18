import { useEventByYear } from "@/lib/api/event";
import { useYear } from "@/lib/hooks/useYear";
import { Event } from "@/lib/types/event";
import { Indeterminate } from "../atoms/Indeterminate";
import { motion } from "framer-motion";
import { EventCard } from "./EventCard";
import { HeadlessCard } from "../molecules/HeadlessCard";
import { CardContent, CardHeader, CardTitle } from "../ui/card";
import { useMemo } from "react";

interface Props {
  selected: number[];
  setSelected: (events: number[]) => void;
  future?: boolean;
}

export function EventSelector({ selected, setSelected, future = false }: Props) {
  const { year } = useYear()
  const { data: eventsAll, isLoading: isLoadingEvents } = useEventByYear(year)

  const now = new Date()
  const events = useMemo(() => eventsAll?.filter(e => !future || e.startTime.getTime() > now.getTime()), [eventsAll])

  if (isLoadingEvents) {
    return <Indeterminate />
  }

  if (!events?.length) {
    return (
      <div className="flex flex-col">
        <span>No events this year</span>
        <span>If this is wrong try manually running the 'update events' task</span>
      </div>
    )
  }

  const handleToggleEvent = (event: Event) => {
    let newSelected

    if (selected.some(id => id === event.id)) {
      newSelected = selected.filter(id => id !== event.id)
    } else {
      newSelected = [...selected]
      newSelected.push(event.id)
    }

    setSelected(newSelected)
  }

  return (
    <HeadlessCard>
      <CardHeader className="px-4 sm:px-0 pt-0">
        <CardTitle>Select events</CardTitle>
      </CardHeader>
      <CardContent className="px-4 sm:px-0">
        <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4 perspective-distant">
          {events?.map(e => {
            const isSelected = selected.some(id => id === e.id)

            return (
              <motion.div
                key={e.id}
                layout
                animate={{ rotateX: isSelected ? 360 : 0 }}
                transition={{ duration: 1 }}
                className="cursor-pointer transform-3d"
              >
                <EventCard onClick={() => handleToggleEvent(e)} event={e} className={isSelected ? "h-full border-primary" : "h-full"} />
              </motion.div>
            )
          })}
        </div>
      </CardContent>
    </HeadlessCard>
  )
}
