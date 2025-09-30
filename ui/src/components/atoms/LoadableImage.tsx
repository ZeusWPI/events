import { ComponentProps, useState } from "react"
import { Skeleton } from "../ui/skeleton";
import { cn } from "@/lib/utils/utils";

export const LoadableImage = ({ className, ...props }: ComponentProps<"img">) => {
  const [loaded, setLoaded] = useState(false);

  return (
    <>
      {!loaded && <Skeleton className={cn("w-full h-full rounded-lg", className)} />}
      <img className={!loaded ? "invisible" : cn("opacity-100", className)} onLoad={() => setLoaded(true)} {...props} />
    </>
  )
}
