import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { EventSelector } from "@/components/events/EventSelector";
import { PageHeader } from "@/components/molecules/PageHeader";
import { Button } from "@/components/ui/button";
import { useEventByYear } from "@/lib/api/event";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { useYear } from "@/lib/hooks/useYear";
import { generatePptx } from "@/lib/pptx/pptx";
import { weightCategory } from "@/lib/types/general";
import { Loader2Icon } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

export function PowerPoints() {
  const { year } = useYear()
  const { data: yearEvents, isLoading: isLoadingEvents } = useEventByYear(year)

  const [eventIds, setEventIds] = useState<number[]>([])
  const [generating, setGenerating] = useState(false)

  useBreadcrumb({ title: "Powerpoints", weight: weightCategory, link: { to: "/powerpoints" } })

  if (isLoadingEvents) {
    return <Indeterminate />
  }

  const handleGenerate = () => {
    setGenerating(true)

    const events = yearEvents?.filter(e => eventIds.includes(e.id)) ?? []

    generatePptx(events)
      .then(() => {
        setEventIds([])
        toast.success("Success")
      })
      .catch((err) => toast.error("Failed", { description: err.message }))
      .finally(() => setGenerating(false))
  }

  return (
    <div className="flex flex-col gap-8">
      <PageHeader>
        <Title>Powerpoints</Title>
        <Button onClick={handleGenerate} size="lg" variant="outline" disabled={generating}>
          {generating
            ? <Loader2Icon className="animate-spin" />
            : <span>Generate</span>
          }
        </Button>
      </PageHeader>
      <EventSelector selected={eventIds} setSelected={setEventIds} />
    </div>
  )
}
