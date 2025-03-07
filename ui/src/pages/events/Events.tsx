import { LoadingCards } from "@/components/organisms/LoadingCards";
import { Button } from "@/components/ui/button";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Navigate, Outlet, useParams } from "@tanstack/react-router";
import { CalendarPlus2 } from "lucide-react";
import Error404 from "../404";

function Events() {
  const yearString = useParams({ from: "/events/$year", shouldThrow: false });

  const { data: years } = useYearGetAll();

  useBreadcrumb({ title: "Events", link: { to: "/events" } });

  if (!years) {
    return (
      <LoadingCards rows={3} cols={4} />
    );
  }

  if (!years.length) {
    return (
      <div className="flex flex-col justify-center items-center h-full space-y-4">
        <CalendarPlus2 className="size-12 text-muted-foreground" />
        <h3 className="text-lg font-semibold">No events found</h3>
        <h5 className="text-md text-muted-foreground">Get started by scraping the website for events</h5>
        <Button>TODO</Button>
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
