import { useAuth } from "@/lib/hooks/useAuth";
import { Forbidden } from "@/pages/Forbidden";
import { Login } from "@/pages/Login";

export function AuthLayout({ children }: { children: React.ReactNode }) {
  const { user, isLoading, forbidden } = useAuth();

  if (isLoading) {
    // Avoid a brief flickering of the login view when you're already logged in
    return null;
  }

  if (forbidden) {
    return <Forbidden />;
  }

  if (!user) {
    return <Login />;
  }

  return children;
}
