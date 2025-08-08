import { Indeterminate } from "@/components/atoms/Indeterminate";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { useParams, Navigate, Outlet } from "@tanstack/react-router";
import Error404 from "../404";
import { weightCategory } from "@/lib/types/general";

export function Announcements() {
  const yearString = useParams({ from: "/announcements/$year", shouldThrow: false });

  const { data: years, isLoading } = useYearGetAll();

  useBreadcrumb({ title: "Announcements", weight: weightCategory });

  if (isLoading) {
    return <Indeterminate />;
  }

  if (!years) {
    return null; // Caught by error component
  }

  const year = years?.find(({ formatted }) => formatted === yearString?.year);
  if (yearString && !year) {
    return <Error404 />;
  }

  return (
    <>
      {!yearString && <Navigate to="/announcements/$year" params={{ year: years[0]!.formatted }} />}
      <Outlet />
    </>
  );
}
