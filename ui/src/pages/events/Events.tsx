import { Navigate, Outlet, useParams } from "@tanstack/react-router";
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
