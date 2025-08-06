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

  useBreadcrumb({ title: "Powerpoints", link: { to: "/powerpoints" } })

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
    setYear(years?.find(y => y.id === Number(value)))
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
          <CardTitle>Selected events</CardTitle>
        </CardHeader>
        <CardContent className="px-0 space-y-4">
          {events.length > 0 ? (
            <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
              {events.map(e => (
                <motion.div
                  key={e.id}
                  layout
                >
                  <EventCard event={e} onClick={() => handleRemoveEvent(e)} />
                </motion.div>
              ))}
            </div>
          ) : (
            <span>Select the events to generate a powerpoint for</span>
          )}
        </CardContent>
      </HeadlessCard>
      <HeadlessCard>
        <CardHeader className="px-4 sm:px-0 pt-0">
          <CardTitle>Add events</CardTitle>
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
              <div className="grid gap-4 lg:grid-cols-2 xl:grid-cols-4">
                {yearEvents
                  ?.filter(e => !events.map(e => e.id).includes(e.id))
                  .map(e => (
                    <motion.div
                      key={e.id}
                      layout
                    >
                      <EventCard event={e} onClick={() => handleAddEvent(e)} />
                    </motion.div>
                  ))}
              </div>
              {yearEvents?.filter(e => !events.map(e => e.id).includes(e.id)).length === 0 && (
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
