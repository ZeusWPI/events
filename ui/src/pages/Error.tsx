import { ErrorComponentProps, Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { isResponseNot200Error } from "@/lib/api/query";
import Error404 from "./404";
import { Forbidden } from "./Forbidden";

export function Error({ error }: ErrorComponentProps) {
  if (isResponseNot200Error(error)) {
    switch (error.response.status) {
      case 404:
        return (
          <div className="w-full h-full pt-[10%]">
            <Error404 />
          </div>
        )
      case 403:
        return (
          <div className="w-full h-full pt-[10%]">
            <Forbidden />
          </div>
        )
    }
  }

  return (
    <div className="flex flex-col justify-center items-center h-full pt-[10%]">
      <p className="font-semibold text-primary">
        500
      </p>
      <h1 className="mt-4 text-balance text-5xl font-semibold tracking-tight sm:text-7xl">
        Server Error
      </h1>
      <p className="mt-6 text-pretty text-lg font-medium text-gray-500 sm:text-xl/8">
        Kapot
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
