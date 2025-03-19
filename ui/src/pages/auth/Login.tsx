import { DividerText } from "@/components/atoms/DividerText";
import EventsIcon from "@/components/icons/EventsIcon";
import GithubIcon from "@/components/icons/GithubIcon";
import WebsiteIcon from "@/components/icons/WebsiteIcon";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/hooks/useAuth";

export function Login() {
  const { login } = useAuth();

  return (
    <div className="flex flex-col items-center justify-center min-h-svh">
      <div className="flex flex-col items-center justify-center w-full max-w-sm gap-4">
        <EventsIcon className="h-8 w-8 fill-secondary-foreground" />
        <span className="text-xl font-bold">Events</span>
        <div className=" flex flex-col gap-2 w-full">
          <Button onClick={login} className="w-full">Login with zauth</Button>
          <DividerText className="py-1">or</DividerText>
          <div className="grid gap-4 sm:grid-cols-2">
            <Button variant="outline" className="w-full" asChild>
              <a href="https://zeus.gent" target="_blank" rel="noopener noreferrer">
                <GithubIcon className="stroke-secondary-foreground" />
                Github
              </a>
            </Button>
            <Button variant="outline" className="w-full" asChild>
              <a href="https://zeus.gent" target="_blank" rel="noopener noreferrer">
                Website
                <WebsiteIcon className="fill-secondary-foreground" />
              </a>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
