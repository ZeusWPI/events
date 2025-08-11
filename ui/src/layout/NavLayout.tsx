import type { ReactNode } from "react";
import AppSidebar from "@/components/organisms/navigation/AppSidebar";
import { useYear } from "@/lib/hooks/useYear";
import { Indeterminate } from "@/components/atoms/Indeterminate";

function NavLayout({ children }: { children: ReactNode }) {
  const { isLoading } = useYear()

  if (isLoading) {
    return <Indeterminate />
  }

  return (
    <AppSidebar>
      {children}
    </AppSidebar>
  );
}

export default NavLayout;
