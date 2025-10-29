import { useEventByYear } from "@/lib/api/event";
import { useOrganizerByYear } from "@/lib/api/organizer";
import { useAuth } from "@/lib/hooks/useAuth";
import { useYear } from "@/lib/hooks/useYear";
import { Event } from "@/lib/types/event";
import { mailSchema, MailSchema } from "@/lib/types/mail";
import { useForm } from "@tanstack/react-form";
import { Link } from "@tanstack/react-router";
import { format } from "date-fns";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { ButtonGroup } from "../atoms/ButtonGroup";
import { FormField } from "../atoms/FormField";
import { Indeterminate } from "../atoms/Indeterminate";
import { Title } from "../atoms/Title";
import { EventSelector } from "../events/EventSelector";
import { HeadlessCard } from "../molecules/HeadlessCard";
import { PageHeader } from "../molecules/PageHeader";
import { DateTimePicker } from "../organisms/DateTimePicker";
import { MarkdownCombo } from "../organisms/markdown/MarkdownCombo";
import { Button } from "../ui/button";
import { CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { nlBE } from "date-fns/locale";
import { capitalize } from "@/lib/utils/utils";

interface Props {
  mail?: MailSchema;
  onSubmit: (mail: MailSchema) => void;
}

export function MailForm({ mail, onSubmit }: Props) {
  const { user } = useAuth()
  const { year } = useYear()
  const { data: events, isLoading: isLoadingEvents } = useEventByYear(year)
  const { data: organizers, isLoading: isLoadingOrganizers } = useOrganizerByYear(year)

  const [referenceDate, setReferenceDate] = useState<Date | undefined>(undefined)


  const handleSubmit = () => {
    const selected = events?.filter(e => form.state.values.eventIds.includes(e.id)) ?? []

    if (selected.some(e => e.startTime.getTime() < form.state.values.sendTime.getTime())) {
      toast.error("Invalid send time", { description: "Mail send time needs to be before every selected event" })
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

  const handleGenerator = (name: string) => {
    const generator = generators.find(g => g.name === name)
    if (!generator) return

    const selectedEventIds = form.getFieldValue("eventIds") as number[]
    const selectedEvents = events?.filter(e => selectedEventIds.includes(e.id)) ?? []
    const text = generator.func(selectedEvents)

    form.setFieldValue("content", text)
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
        <HeadlessCard>
          <CardHeader className="px-4 sm:px-0 pt-0">
            <CardTitle>Generate</CardTitle>
            <CardDescription className="flex flex-col">
              This will generate some text based on the selected events
              <span className="font-bold">It will overwrite all current content!</span>
            </CardDescription>
          </CardHeader>
          <CardContent className="px-4 sm:px-0">
            <ButtonGroup>
              {generators.map(g => (
                <Button key={g.name} type="button" onClick={() => handleGenerator(g.name)} variant="outline" disabled={!(form.getFieldValue("eventIds") as number[]).length}>{g.name}</Button>
              ))}
            </ButtonGroup>
          </CardContent>
        </HeadlessCard>
        <form.Field name="eventIds" listeners={{ onChange: ({ value }) => updateReferenceDate(value as number[]) }}>
          {(field) => (
            <EventSelector selected={field.state.value as number[]} setSelected={field.handleChange} />
          )}
        </form.Field>
      </form>
    </div>
  )
}

type Generator = {
  name: string;
  func: (events: Event[]) => string;
}

const generators: Generator[] = [
  {
    name: "Short",
    func: (events: Event[]) => events.map(e => capitalize(`${format(e.startTime, "iiii (dd LLLL)", { locale: nlBE })}: **[${e.name}](${e.url})**`)).join("\n\n")
  },
  {
    name: "Long",
    func: (events: Event[]) => events.map(e => `### ${e.name}\n\n__🕑 ${capitalize(format(e.startTime, "iiii (dd LLLL)", { locale: nlBE }))}__ \\\n__📍 ${e.location}__`).join("\n\n\n\n")
  },
]
