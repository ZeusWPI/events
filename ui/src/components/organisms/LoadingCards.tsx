import clsx from "clsx";
import { SkeletonCard } from "../molecules/SkeletonCard";

interface Props {
  rows?: number;
  cols?: number;
}

export function LoadingCards({ rows = 1, cols = 1 }: Props) {
  return (
    <div
      className={clsx(
        "grid gap-10",
        `grid-rows-${rows}`,
        `grid-cols-${cols}`,
      )}
    >
      {Array.from({ length: rows * cols }).map((_, i) => (
        <SkeletonCard key={i} /> // eslint-disable-line react/no-array-index-key
      ))}
    </div>
  );
}
