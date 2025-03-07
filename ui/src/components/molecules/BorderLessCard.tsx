import type { ReactNode } from "react";
import { Card } from "../ui/card";

export function BorderlessCard({ children }: { children: ReactNode }) {
  return (
    <Card className="border-none w-full">
      {children}
    </Card>
  );
}
