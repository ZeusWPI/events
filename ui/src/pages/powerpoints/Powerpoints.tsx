import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { EventCard } from "@/components/events/EventCard";
import { HeadlessCard } from "@/components/molecules/HeadlessCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useEventByYear } from "@/lib/api/event";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { generatePptx } from "@/lib/pptx/pptx";
import { Event } from "@/lib/types/event";
import { weightCategory } from "@/lib/types/general";
import { motion } from "framer-motion";
import { Loader2Icon } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

export function PowerPoints() {
  const { data: years, isLoading: isLoadingYears } = useYearGetAll()
  const [year, setYear] = useState(years?.[0])
  const { data: yearEvents, isLoading: isLoadingEvents } = useEventByYear({ id: year?.id ?? 0 })

  const [events, setEvents] = useState<Event[]>([])
  const [generating, setGenerating] = useState(false)

  useBreadcrumb({ title: "Powerpoints", weight: weightCategory, link: { to: "/powerpoints" } })

  if (isLoadingYears) {
    return <Indeterminate />
  }

  const handleGenerate = () => {
    setGenerating(true)
    generatePptx(events)
      .then(() => {
        setEvents([])
        toast.success("Success")
      })
      .catch((err) => toast.error("Failed", { description: err.message }))
      .finally(() => setGenerating(false))
  }

  const handleSelectChange = (value: string) => {
    const newYear = years?.find(y => y.id === Number(value))
    if (newYear?.id === year?.id) {
      return
    }

    setEvents([])
    setYear(newYear)
  }

  const handleToggleEvent = (event: Event) => {
    if (events.some(e => e.id === event.id)) {
      handleRemoveEvent(event)
    } else {
      handleAddEvent(event)
    }
  }

  const handleAddEvent = (event: Event) => {
    const newEvents = [...events]
    newEvents.push(event)
    newEvents.sort((a, b) => a.startTime.getTime() - b.startTime.getTime())
    setEvents(newEvents)
  }

  const handleRemoveEvent = (event: Event) => {
    const newEvents = events.filter(e => e.id !== event.id)
    setEvents(newEvents)
  }

  return (
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>Powerpoints</Title>
        <Button onClick={handleGenerate} size="lg" variant="outline" disabled={generating || !events.length}>
          {generating
            ? <Loader2Icon className="animate-spin" />
            : <span>Generate</span>
          }
        </Button>
      </PageHeader>
      <HeadlessCard>
        <CardHeader className="px-4 sm:px-0 pt-0">
          <CardTitle>Select events</CardTitle>
        </CardHeader>
        <CardContent className="px-0 space-y-4">
          <Select onValueChange={handleSelectChange} defaultValue={year?.id.toString()}>
            <SelectTrigger>
              <SelectValue placeholder="Select year" />
            </SelectTrigger>
            <SelectContent>
              {years?.map(y => (
                <SelectItem key={y.id} value={y.id.toString()}>
                  {y?.formatted}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {isLoadingEvents ? (
            <Indeterminate />
          ) : (
            <>
              <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4 perspective-distant">
                {yearEvents?.map(e => {
                  const selected = events.some(ev => ev.id === e.id)

                  return (
                    <motion.div
                      key={e.id}
                      layout
                      onClick={() => handleToggleEvent(e)}
                      animate={{ rotateX: selected ? 360 : 0 }}
                      transition={{ duration: 1 }}
                      className="cursor-pointer"
                    >
                      <EventCard event={e} onClick={() => handleToggleEvent(e)} className={selected ? "h-full border-primary" : "h-full"} />
                    </motion.div>
                  )
                })}
              </div>
              {yearEvents?.length === 0 && (
                <div className="flex flex-col">
                  <span>No event this year</span>
                  <span>Please select a different year</span>
                </div>
              )}
            </>
          )}
        </CardContent>
      </HeadlessCard>
    </div>
  )
}
