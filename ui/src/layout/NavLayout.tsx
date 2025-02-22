import type { ReactNode } from "react";
import AppSidebar from "@/components/organisms/navigation/AppSidebar";

function NavLayout({ children }: { children: ReactNode }) {
  return (
    <AppSidebar>
      {children}
    </AppSidebar>
  );
}

export default NavLayout;
