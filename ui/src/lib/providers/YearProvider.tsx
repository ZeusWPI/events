import type { ReactNode } from "react";
import { useEffect, useMemo, useState } from "react";
import { useYearGetAll } from "../api/year";
import { YearContext } from "../contexts/yearContext";
import { Year } from "../types/year";
import { toast } from "sonner";
import { useAuth } from "../hooks/useAuth";

interface YearProviderProps {
  children: ReactNode;
  storageKey?: string;
}

export function YearProvider({ children, storageKey = "events-ui-year" }: YearProviderProps) {
  const { user } = useAuth()
  const [year, setYear] = useState<Year>({ id: 0, start: 0, end: 0, formatted: "" })
  const [locked, setLocked] = useState(false)

  const { data, isLoading } = useYearGetAll(!!user);

  useEffect(() => {
    if (data && data.length) {
      const cookieYearId = Number(localStorage.getItem(storageKey))
      const cookieYear = data.find(y => y.id === cookieYearId)

      if (cookieYear) {
        setYear(cookieYear)
      } else {
        setYear(data[0]!)
      }
    }
  }, [data, storageKey])

  const value = useMemo(
    () => ({
      year,
      setYear: (year: Year) => {
        if (locked) {
          toast.error("Year locked", { description: "Year cannot be changed on this page" })
          return
        }

        localStorage.setItem(storageKey, year.id.toString())
        setYear(year)
      },
      isLoading,
      locked,
      setLocked,
    }),
    [year, isLoading, locked, storageKey]
  )

  return <YearContext value={value}>{children}</YearContext>;
}
