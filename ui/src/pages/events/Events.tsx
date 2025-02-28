import { LoadingCards } from "@/components/molecules/LoadingCards";
import { Button } from "@/components/ui/button";
import { useYearGetAll } from "@/lib/api/year";
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb";
import { Navigate, Outlet, useParams } from "@tanstack/react-router";
import { CalendarPlus2 } from "lucide-react";
import Error404 from "../404";

function Events() {
  const yearId = useParams({ from: "/events/$yearId", shouldThrow: false });

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

  const year = years?.find(({ id }) => id.toString() === yearId?.yearId);
  if (yearId && !year) {
    return <Error404 />;
  }

  return (
    <>
      {!yearId && <Navigate to="/events/$yearId" params={{ yearId: years[0]!.id.toString() }} />}
      <Outlet />
    </>
  );
}

export default Events;
