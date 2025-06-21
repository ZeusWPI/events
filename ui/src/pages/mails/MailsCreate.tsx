import { Indeterminate } from "@/components/atoms/Indeterminate";
import { Title } from "@/components/atoms/Title";
import { EventCard } from "@/components/events/EventCard";
import { HeadlessCard } from "@/components/molecules/HeadlessCard";
import { PageHeader } from "@/components/molecules/PageHeader";
import { DateTimePicker } from "@/components/organisms/DateTimePicker";
import { MarkdownCombo } from "@/components/organisms/markdown/MarkdownCombo";
import { Button } from "@/components/ui/button";
import { CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useEventByYear } from "@/lib/api/event";
import { useMailCreate, useMailGetAll, useMailUpdate } from "@/lib/api/mail";
import { useYearGetAll } from "@/lib/api/year";
import { Event } from "@/lib/types/event";
import { Link, useNavigate, useParams } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { useState } from "react";
import { toast } from "sonner";

export function MailsCreate() {
  const params = useParams({ from: "/mails/edit/$mail", shouldThrow: false })
  const mailId = params?.mail ? Number(params.mail) : 0

  const { data: mails, isLoading: isLoadingMails } = useMailGetAll()
  const mail = mails?.find(m => m.id === mailId)

  const mailCreate = useMailCreate()
  const mailUpdate = useMailUpdate()

  const { data: years, isLoading: isLoadingYears } = useYearGetAll()
  const [year, setYear] = useState(years?.[0])
  const { data: yearEvents, isLoading: isLoadingEvents } = useEventByYear({ id: year?.id ?? 0 })

  const [submitting, setSubmitting] = useState(false)
  const [date, setDate] = useState(mail?.sendTime)
  const [title, setTitle] = useState(mail?.title ?? "")
  const [content, setContent] = useState(mail?.content)
  const [events, setEvents] = useState(mail?.events ?? [])

  const navigate = useNavigate()

  if (isLoadingMails || isLoadingYears) {
    return <Indeterminate />
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

  const handleSubmit = () => {
    const now = Date.now()

    if (!date || date.getTime() <= now) {
      toast.error("Invalid date", { description: "Date has to be in the future" })
      return
    }

    if (!title || !title.length) {
      toast.error("Invalid title", { description: "Title can't be empty" })
      return
    }

    if (!content || !content.length) {
      toast.error("Invalid mail", { description: "Mail text can't be empty" })
      return
    }

    if (!events.length) {
      toast.error("Invalid events", { description: "Select at least one event" })
      return
    }

    setSubmitting(true)

    let action
    if (mailId) {
      action = mailUpdate
    } else {
      action = mailCreate
    }

    action.mutate({ mail: { id: mailId, title, content, sendTime: date }, eventIds: events.map(e => e.id) }, {
      onSuccess: () => {
        toast.success("Success")
        navigate({ to: "/mails" })
      },
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setSubmitting(false),
    })
  }

  return (
    <div className="space-y-8">
      <PageHeader>
        <Title>{`${mailId ? "Edit" : "Create"} Mail`}</Title>
        <div className="flex justify-center gap-2">
          <Button variant="outline" asChild>
            <Link to="/mails">
              Cancel
            </Link>
          </Button>
          <Button onClick={handleSubmit} disabled={submitting || (date?.getTime() === mail?.sendTime.getTime() && content === mail?.content && events.map(e => e.id) === mail?.events.map(e => e.id))}>
            Submit
          </Button>
        </div>
      </PageHeader>
      <HeadlessCard>
        <CardHeader className="px-4 sm:px-0 pt-0">
          <CardTitle>Mail</CardTitle>
        </CardHeader>
        <CardContent className="px-0 space-y-4">
          <div className="grid grid-cols-[auto_1fr] items-center gap-2 space-x-4">
            <span>Send time</span>
            <DateTimePicker value={date} onChange={setDate} weekStartsOn={1} />
            <span>Title</span>
            <Input placeholder="[Zeus WPI] Events update" onChange={event => setTitle(event.target.value)} />
          </div>
          <MarkdownCombo value={content} onChange={setContent} textareaProps={{ placeholder: "Write mail here..." }} />
        </CardContent>
      </HeadlessCard>
      <HeadlessCard>
        <CardHeader className="px-4 sm:px-0 pt-0">
          <CardTitle>Covered events</CardTitle>
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
            <span>Select the events that this mail covers</span>
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
          )}
        </CardContent>
      </HeadlessCard>
    </div>
  )
}
