import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";

const Error404: React.FC = () => {
  return (
    <div className="flex flex-col justify-center items-center h-full">
      <p className="font-semibold text-primary">
        404
      </p>
      <h1 className="mt-4 text-balance text-5xl font-semibold tracking-tight sm:text-7xl">
        Page not found
      </h1>
      <p className="mt-6 text-pretty text-lg font-medium text-gray-500 sm:text-xl/8">
        Sorry, the page you requested couldn't be found.
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
};

export default Error404;
