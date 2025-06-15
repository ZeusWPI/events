import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";

export function Error() {
  return (
    <div className="flex flex-col justify-center items-center h-full pt-[10%]">
      <p className="font-semibold text-primary">
        500
      </p>
      <h1 className="mt-4 text-balance text-5xl font-semibold tracking-tight sm:text-7xl">
        Server Error
      </h1>
      <p className="mt-6 text-pretty text-lg font-medium text-gray-500 sm:text-xl/8">
        Sorry, developer skill issue
      </p>
      <div className="mt-10 flex items-center justify-center gap-x-6">
        <Button asChild>
          <Link to="/">
            Go back home
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
