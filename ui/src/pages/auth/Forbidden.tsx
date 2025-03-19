import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/hooks/useAuth";
import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";

export function Forbidden() {
  const { logout } = useAuth();

  return (
    <div className="flex flex-col justify-center items-center h-full pt-[10%]">
      <p className="font-semibold text-primary">
        403
      </p>
      <h1 className="mt-4 text-balance text-5xl font-semibold tracking-tight sm:text-7xl">
        Forbidden
      </h1>
      <p className="mt-6 text-pretty text-lg font-medium text-gray-500 sm:text-xl/8">
        You don't have the required permissions to view this website.
        <br />
        Contact the board if you think this is a mistake.
      </p>
      <div className="mt-10 flex items-center justify-center gap-x-6">
        <Button onClick={logout} asChild>
          <Link to="/">
            Try again
          </Link>
        </Button>
        <Button asChild variant="ghost">
          <a href="https://zeus.gent" target="_blank" rel="noopener noreferrer">
            Zeus WPI website
            <ArrowRight />
          </a>
        </Button>
      </div>
    </div>
  );
}
