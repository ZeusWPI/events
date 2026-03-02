import { useEventByYear } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useAuth } from "@/lib/hooks/useAuth";
import { useYear } from "@/lib/hooks/useYear";
import { announcementSchema, AnnouncementSchema } from "@/lib/types/announcement";
import { useForm, useStore } from "@tanstack/react-form";
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
import { Confirm } from "../molecules/Confirm";
import { Checkbox } from "../ui/checkbox";
import { Label } from "../ui/label";

interface Props {
  announcement?: AnnouncementSchema;
  defaultEvents?: number[];
  onSubmit: (announcement: AnnouncementSchema) => void;
}

// Returns true if the text contains any of the mentions
const includesMention = (text: string, mentions: string[]): boolean => {
  for (const m of mentions) {
    if (text.includes(m)) return true
  }

  return false
}

export function AnnouncementForm({ announcement, defaultEvents, onSubmit }: Props) {
  const { user } = useAuth()
  const { year } = useYear()
  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year)
  const { data: organizers, isLoading: isLoadingOrganizers } = useOrganizerByYear(year)

  const [referenceDate, setReferenceDate] = useState<Date | undefined>(undefined)

  const [openAtChannel, setOpenAtChannel] = useState(false)
  const handleAtChannel = () => {
    handleSubmit(true)
  }

  const handleSubmit = (ignoreAtChannel: boolean) => {
    if (form.state.values.draft) {
      form.setFieldValue("sendTime", undefined)
    } else {
      const selected = events?.filter(e => form.state.values.eventIds.includes(e.id)) ?? []

      if (selected.some(e => e.startTime.getTime() < (form.state.values.sendTime?.getTime() ?? 0))) {
        toast.error("Invalid send time", { description: "Announcement send time needs to be before every selected event" })
        return
      }
    }

    const organizer = organizers?.find(o => o.id === user?.id)
    if (!organizer) {
      toast.error("Not a board member", { description: "You were not a board member that year" })
      return
    }

    if (!ignoreAtChannel) {
      // Check mattermost mention
      let mentions = ["@channel"]
      if (!includesMention(form.state.values.content, mentions)) {
        setOpenAtChannel(true)
        return
      }
      // Check discord mentions
      mentions = ["@here", "@everyone"]
      if (!includesMention(form.state.values.content, mentions)) {
        setOpenAtChannel(true)
        return
      }
    }

    onSubmit(form.state.values)
  }

  const form = useForm({
    defaultValues: announcement ?? {
      yearId: year.id,
      eventIds: defaultEvents ?? [],
      content: "",
      sendTime: new Date(),
      draft: false,
    },
    validators: {
      onSubmit: announcementSchema,
    },
    onSubmit: () => handleSubmit(false),
  })

  const isDraft = useStore(form.store, (s) => s.values.draft)

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
    <>
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
        <div className="border border-orange-500 rounded-lg whitespace-pre p-8">
          {`Mattermost supports only a small subset of markdown features.\nThe preview you see does not necessarily reflect what mattermost will render.\nYou can find `}
          {`a list of supported features `}
          <a href="https://mattermost.com/blog/laymans-guide-to-markdown-on-mattermost/" target="_blank" rel="noopener noreferrer" className="cursor-pointer underline underline-offset-4 decoration-orange-500 hover:no-underline">here</a>
          {`.`}
        </div>
        <form className="space-y-4" onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}>
          <div className="grid grid-cols-[auto_1fr] items-center gap-2 space-x-4">
            <Label htmlFor="announcement-form-draft">Draft</Label>
            <form.Field name="draft">
              {(field) => (
                <FormField field={field}>
                  <Checkbox id="announcement-form-draft" checked={field.state.value as boolean} onCheckedChange={field.handleChange} />
                </FormField>
              )}
            </form.Field>
            <Label htmlFor="announcement-form-send-time">Send time</Label>
            <form.Field name="sendTime">
              {(field) => (
                <FormField field={field} className="flex items-center gap-4">
                  <DateTimePicker id="announcement-form-send-time" value={field.state.value as Date} setValue={field.handleChange} referenceDate={referenceDate} inputProps={{ disabled: isDraft }} />
                </FormField>
              )}
            </form.Field>
          </div>
          <form.Field name="content">
            {(field) => (
              <FormField field={field}>
                <MarkdownCombo value={field.state.value as string} onChange={field.handleChange} textareaProps={{ placeholder: "Write announcement here..." }} />
              </FormField>
            )}
          </form.Field>
          <form.Field name="eventIds" listeners={{ onChange: ({ value }) => updateReferenceDate(value as number[]) }}>
            {(field) => (
              <EventSelector selected={field.state.value as number[]} setSelected={field.handleChange} future />
            )}
          </form.Field>
        </form>
      </div>
      <Confirm
        title="Send confirmation"
        description="The announcement has no mattermost or discord mentions. Do you still want to send it?"
        confirmText="Send"
        onConfirm={handleAtChannel}
        open={openAtChannel}
        onOpenChange={setOpenAtChannel}
      />
    </>
  )
}
