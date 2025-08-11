import { formatDate, formatTime } from "@/lib/utils/utils";
import { CalendarIcon } from "lucide-react";
import { ChangeEvent, useMemo, useState } from "react";
import { Button } from "../ui/button";
import { Calendar } from "../ui/calendar";
import { Card, CardContent, CardFooter } from "../ui/card";
import { Input } from "../ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { parse } from "date-fns";

interface Props {
  value: Date;
  setValue: (date: Date) => void;
  id?: string;
  referenceDate?: Date;
}

type Preset = {
  label: string;
  date: Date;
}

function getNextWeekdayTime(start: Date, weekday: number, hours: number, minutes = 0) {
  const date = new Date(start.getTime());
  date.setHours(hours, minutes, 0, 0);
  const diff = (weekday - date.getDay() + 7) % 7 || 7;
  date.setDate(date.getDate() + diff);

  return date;
}

function getPrevWeekdayTime(start: Date, weekday: number, hours: number, minutes = 0) {
  const date = new Date(start.getTime());
  date.setHours(hours, minutes, 0, 0);
  const diff = (date.getDay() - weekday + 7) % 7 || 7;
  date.setDate(date.getDate() - diff);

  return date;
}

function getPresets(referenceDate?: Date): Preset[] {
  const presets: Preset[] = []
  const now = new Date()

  presets.push({
    label: "Next Saturday",
    date: getNextWeekdayTime(now, 6, 20),
  });
  presets.push({
    label: "Next Sunday",
    date: getNextWeekdayTime(now, 0, 20),
  });

  if (referenceDate) {
    presets.push({
      label: "Event Sunday",
      date: getPrevWeekdayTime(referenceDate, 0, 20),
    });
    presets.push({
      label: "Event Saturday",
      date: getPrevWeekdayTime(referenceDate, 6, 20),
    });
  }

  return presets
}

export function DateTimePicker({ value, setValue, id, referenceDate }: Props) {
  const [open, setOpen] = useState(false)

  const presets = useMemo(() => getPresets(referenceDate), [referenceDate])

  const handleSetDate = (date: Date) => {
    const newDate = new Date(value.getTime())
    newDate.setFullYear(date.getFullYear(), date.getMonth(), date.getDate())
    setValue(newDate)
  }

  const handleSetTime = (e: ChangeEvent<HTMLInputElement>) => {
    const newDate = new Date(value?.getTime())
    const time = parse(e.target.value, "HH:mm:ss", new Date())

    newDate.setHours(time.getHours())
    newDate.setMinutes(time.getMinutes())
    newDate.setSeconds(time.getSeconds())

    setValue(newDate)
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div className="relative flex gap-2">
          <Input id={id} className="bg-background pr-10" value={formatDate(value)} readOnly />
          <Button type="button" size="icon" variant="ghost" className="absolute top-1/2 right-2 size-6 -translate-y-1/2">
            <CalendarIcon className="size-3.5" />
          </Button>
        </div>
      </PopoverTrigger>
      <PopoverContent align="end" alignOffset={-8} sideOffset={10} className="w-auto overflow-hidden p-0">
        <Card className="w-fit py-3">
          <CardContent className="relative p-0 md:pr-48">
            <Calendar mode="single" selected={value} onSelect={handleSetDate} required className="bg-transparent p-0 [--cell-size:--spacing(10.5)]" />
            <div className="no-scrollbar inset-y-0 right-0 flex max-h-72 w-full scroll-pb-6 flex-col gap-4 overflow-y-auto border-t p-6 md:absolute md:max-h-none md:w-48 md:border-t-0 md:border-l">
              <div className="grid gap-2">
                {presets.map((preset) => (
                  <Button
                    key={preset.label}
                    type="button"
                    variant="outline"
                    size="sm"
                    className="flex-1"
                    onClick={() => setValue(preset.date)}
                  >
                    {preset.label}
                  </Button>
                ))}
              </div>
            </div>
          </CardContent>
          <CardFooter className="border-t">
            <Input type="time" step="1" value={formatTime(value)} onChange={handleSetTime} />
          </CardFooter>
        </Card>
      </PopoverContent>
    </Popover>
  )
}
