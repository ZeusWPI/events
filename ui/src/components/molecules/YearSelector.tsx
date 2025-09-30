import { useYearGetAll } from "@/lib/api/year";
import { useYear } from "@/lib/hooks/useYear";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { useSidebar } from "../ui/sidebar";
import { ComponentProps } from "react";
import { cn } from "@/lib/utils/utils";

export function YearSelector() {
  const { state } = useSidebar()
  const isOpen = state === "expanded"

  const { year, setYear, locked } = useYear()
  const { data: years } = useYearGetAll()
  const isCurrentYear = year.id === (years?.[0]?.id ?? -1)

  const handleSelectChange = (value: string) => {
    const newYear = years?.find(y => y.id === Number(value))
    if (!newYear || newYear?.id === year?.id) {
      return
    }

    setYear(newYear)
  }

  const handleReset = () => {
    if (!years || !years.length) {
      return
    }

    setYear(years[0]!)
  }

  return (
    <div className="flex flex-col gap-2">
      {isOpen && locked && <InfoMessage className="border-muted-foreground">Year locked on this page</InfoMessage>}
      {!isCurrentYear && (isOpen
        ? <InfoMessage onClick={handleReset} className="border-red-500 text-red-500 cursor-pointer">Old academic year</InfoMessage>
        : <InfoMessage onClick={handleReset} className="border-red-500 text-red-500 cursor-pointer">!</InfoMessage>
      )}
      <Select onValueChange={handleSelectChange} value={year?.id.toString()} disabled={locked}>
        <SelectTrigger className="w-full">
          {isOpen && <SelectValue />}
        </SelectTrigger>
        <SelectContent className="max-h-72">
          {years?.map(y => (
            <SelectItem key={y.id} value={y.id.toString()}>
              {y?.formatted}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}

function InfoMessage({ onClick, className, ...props }: ComponentProps<'p'>) {
  return <p onClick={onClick} className={cn("border rounded-md p-1 text-sm text-center", className)} {...props} />
}
