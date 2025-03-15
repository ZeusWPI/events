import { DividerText } from "@/components/atoms/DividerText";
import EventsIcon from "@/components/atoms/icons/EventsIcon";
import GithubIcon from "@/components/atoms/icons/GithubIcon";
import WebsiteIcon from "@/components/atoms/icons/WebsiteIcon";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/hooks/useAuth";

export function Login() {
  const { login } = useAuth();

  const handleGithubClick = () => {
    window.location.replace("https://github.com/ZeusWPI/events");
  };

  const handleWebsiteClick = () => {
    window.location.replace("https://zeus.gent");
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-svh">
      <div className="flex flex-col items-center justify-center w-full max-w-sm gap-4">
        <EventsIcon className="h-8 w-8 fill-secondary-foreground" />
        <span className="text-xl font-bold">Events</span>
        <div className=" flex flex-col gap-2 w-full">
          <Button onClick={login} className="w-full">Login with zauth</Button>
          <DividerText className="py-1">or</DividerText>
          <div className="grid gap-4 sm:grid-cols-2">
            <Button onClick={handleGithubClick} variant="outline" className="w-full">
              <GithubIcon className="stroke-secondary-foreground" />
              Github
            </Button>
            <Button onClick={handleWebsiteClick} variant="outline" className="w-full">
              Website
              <WebsiteIcon className="fill-secondary-foreground" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
