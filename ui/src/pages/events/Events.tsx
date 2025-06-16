import { Navigate, Outlet, useParams } from "@tanstack/react-router";
import { CalendarPlus2 } from "lucide-react";
import { Indeterminate } from "@/components/atoms/Indeterminate";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import Error404 from "../404";

function Events() {
  const yearString = useParams({ from: "/events/$year", shouldThrow: false });

  const { data: years, isLoading } = useYearGetAll();

  useBreadcrumb({ title: "Events", link: { to: "/events" } });

  if (isLoading) {
    return <Indeterminate />;
  }

  if (!years) {
    return null; // Caught by error component
  }

  if (!years.length) {
    return (
      <div className="flex flex-col justify-center items-center h-full space-y-4">
        <CalendarPlus2 className="size-12 text-muted-foreground" />
        <h3 className="text-lg font-semibold">No events found</h3>
        <h5 className="text-md text-muted-foreground">Get started by scraping the website for events</h5>
      </div>
    );
  }

  const year = years?.find(({ formatted }) => formatted === yearString?.year);
  if (yearString && !year) {
    return <Error404 />;
  }

  return (
    <>
      {!yearString && <Navigate to="/events/$year" params={{ year: years[0]!.formatted }} />}
      <Outlet />
    </>
  );
}

export default Events;
