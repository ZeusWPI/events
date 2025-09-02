import { useYear } from "@/lib/hooks/useYear";

export function YearLayout({ children }: { children: React.ReactNode }) {
  const { year } = useYear()

  if (!year) {
    return null
  }

  return children
}
