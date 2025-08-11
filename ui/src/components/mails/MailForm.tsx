import { useEventByYear } from "@/lib/api/event";
import { useYear } from "@/lib/hooks/useYear";
import { mailSchema, MailSchema } from "@/lib/types/mail";
import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import { FormField } from "../atoms/FormField";
import { Indeterminate } from "../atoms/Indeterminate";
import { Title } from "../atoms/Title";
import { EventSelector } from "../events/EventSelector";
import { PageHeader } from "../molecules/PageHeader";
import { DateTimePicker } from "../organisms/DateTimePicker";
import { MarkdownCombo } from "../organisms/markdown/MarkdownCombo";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Link } from "@tanstack/react-router";
import { Label } from "../ui/label";
import { useEffect, useState } from "react";

interface Props {
  mail?: MailSchema;
  onSubmit: (mail: MailSchema) => void;
}

export function MailForm({ mail, onSubmit }: Props) {
  const { year } = useYear()
  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year)

  const [referenceDate, setReferenceDate] = useState<Date | undefined>(undefined)


  const handleSubmit = () => {
    const selected = events?.filter(e => form.state.values.eventIds.includes(e.id)) ?? []

    if (selected.some(e => e.startTime.getTime() < form.state.values.sendTime.getTime())) {
      toast.error("Invalid send time", { description: "Mail send time needs to be before every selected event" })
      return
    }

    onSubmit(form.state.values)
  }

  const form = useForm({
    defaultValues: mail ?? {
      yearId: year.id,
      eventIds: [],
      title: "",
      content: "",
      sendTime: new Date(),
    },
    validators: {
      onSubmit: mailSchema,
    },
    onSubmit: handleSubmit,
  })

  const updateReferenceDate = (eventIds: number[]) => {
    const selected = events?.filter(e => eventIds.includes(e.id)) ?? []
    setReferenceDate(selected.sort((a, b) => a.startTime.getTime() - b.startTime.getTime())[0]?.startTime)
  }

  useEffect(() => {
    if (!mail || !events) return
    updateReferenceDate(mail.eventIds)
  }, [mail, events]) // eslint-disable-line react-hooks/exhaustive-deps

  if (isLoadingEvents) {
    return <Indeterminate />
  }

  return (
    <div className="space-y-8">
      <PageHeader>
        <Title>{`${mail ? "Edit" : "Create"} Mail`}</Title>
        <div className="flex justify-center gap-2">
          <Button variant="outline" asChild>
            <Link to="/mails" >
              Cancel
            </Link>
          </Button>
          <Button onClick={form.handleSubmit}>
            Submit
          </Button>
        </div>
      </PageHeader>
      <form className="space-y-4" onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}>
        <div className="grid grid-cols-[auto_1fr] items-center gap-2 space-x-4">
          <Label htmlFor="mail-form-send-time">Send time</Label>
          <form.Field name="sendTime">
            {(field) => (
              <FormField field={field} className="flex items-center gap-4">
                <DateTimePicker id="mail-form-send-time" value={field.state.value as Date} setValue={field.handleChange} referenceDate={referenceDate} />
              </FormField>
            )}
          </form.Field>
          <Label htmlFor="mail-form-title">Title</Label>
          <form.Field name="title">
            {(field) => (
              <FormField field={field}>
                <Input id="mail-form-title" defaultValue={field.state.value as string} placeholder="[Zeus WPI] Events update" onChange={(event) => field.handleChange(event.target.value)} />
              </FormField>
            )}
          </form.Field>
        </div>
        <form.Field name="content">
          {(field) => (
            <FormField field={field}>
              <MarkdownCombo value={field.state.value as string} onChange={field.handleChange} textareaProps={{ placeholder: "Write mail here..." }} />
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


