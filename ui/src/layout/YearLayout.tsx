import { useYear } from "@/lib/hooks/useYear";

export function YearLayout({ children }: { children: React.ReactNode }) {
  const { year, isLoading } = useYear()

  if (isLoading || !year) {
    return null
  }

  return children
}
