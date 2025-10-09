import { useEventByYear } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useAuth } from "@/lib/hooks/useAuth";
import { useYear } from "@/lib/hooks/useYear";
import { announcementSchema, AnnouncementSchema } from "@/lib/types/announcement";
import { useForm } from "@tanstack/react-form";
import { Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { FormField } from "../atoms/FormField";
import { Indeterminate } from "../atoms/Indeterminate";
import { Title } from "../atoms/Title";
import { EventSelector } from "../events/EventSelector";
import { PageHeader } from "../molecules/PageHeader";
import { DateTimePicker } from "../organisms/DateTimePicker";
import { MarkdownCombo } from "../organisms/markdown/MarkdownCombo";
import { Button } from "../ui/button";

interface Props {
  announcement?: AnnouncementSchema;
  onSubmit: (announcement: AnnouncementSchema) => void;
}

export function AnnouncementForm({ announcement, onSubmit }: Props) {
  const { user } = useAuth()
  const { year } = useYear()
  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year)
  const { data: organizers, isLoading: isLoadingOrganizers } = useOrganizerByYear(year)

  const [referenceDate, setReferenceDate] = useState<Date | undefined>(undefined)

  const handleSubmit = () => {
    const selected = events?.filter(e => form.state.values.eventIds.includes(e.id)) ?? []

    if (selected.some(e => e.startTime.getTime() < form.state.values.sendTime.getTime())) {
      toast.error("Invalid send time", { description: "Announcement send time needs to be before every selected event" })
      return
    }

    const organizer = organizers?.find(o => o.id === user?.id)
    if (!organizer) {
      toast.error("Not a board member", { description: "You were not a board member that year" })
      return
    }

    onSubmit(form.state.values)
  }

  const form = useForm({
    defaultValues: announcement ?? {
      yearId: year.id,
      eventIds: [],
      content: "",
      sendTime: new Date(),
    },
    validators: {
      onSubmit: announcementSchema,
    },
    onSubmit: handleSubmit,
  })

  const updateReferenceDate = (eventIds: number[]) => {
    const selected = events?.filter(e => eventIds.includes(e.id)) ?? []
    setReferenceDate(selected.sort((a, b) => a.startTime.getTime() - b.startTime.getTime())[0]?.startTime)
  }

  useEffect(() => {
    if (!announcement || !events) return
    updateReferenceDate(announcement.eventIds)
  }, [announcement, events]) // eslint-disable-line react-hooks/exhaustive-deps

  if (isLoadingEvents) {
    return <Indeterminate />
  }

  return (
    <div className="space-y-8">
      <PageHeader>
        <Title>{`${announcement ? "Edit" : "Create"} Announcement`}</Title>
        <div className="flex justify-center gap-2">
          <Button variant="outline" asChild>
            <Link to="/announcements" >
              Cancel
            </Link>
          </Button>
          <Button onClick={form.handleSubmit} disabled={isLoadingOrganizers}>
            Submit
          </Button>
        </div>
      </PageHeader>
      <form className="space-y-4" onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}>
        <form.Field name="sendTime">
          {(field) => (
            <FormField field={field} className="flex items-center gap-4">
              <label htmlFor={field.name}>Send time</label>
              <DateTimePicker id={field.name} value={field.state.value as Date} setValue={field.handleChange} referenceDate={referenceDate} />
            </FormField>
          )}
        </form.Field>
        <form.Field name="content">
          {(field) => (
            <FormField field={field}>
              <MarkdownCombo value={field.state.value as string} onChange={field.handleChange} textareaProps={{ placeholder: "Write announcement here..." }} />
            </FormField>
          )}
        </form.Field>
        <form.Field name="eventIds" listeners={{ onChange: ({ value }) => updateReferenceDate(value as number[]) }}>
          {(field) => (
            <EventSelector selected={field.state.value as number[]} setSelected={field.handleChange} />
          )}
        </form.Field>
      </form>
    </div>
  )
}
