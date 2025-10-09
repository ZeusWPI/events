import { PropsWithChildren } from "react"
import { Indeterminate } from "./Indeterminate"

interface Props extends PropsWithChildren {
  loading?: boolean
}

export const LoadingOverlay = ({ loading = false, children }: Props) => {
  return (
    <div className="relative">
      {children}

      {loading && (
        <div className="absolute inset-0 flex items-center justify-center bg-background/20 backdrop-blur-xs">
          <div className="flex flex-col items-center gap-2 text-muted-foreground">
            <div className="w-24">
              <Indeterminate />
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
