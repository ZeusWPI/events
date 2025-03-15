import { useAuth } from "@/lib/hooks/useAuth";
import { Login } from "@/pages/auth/Login";

export function AuthLayout({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    // Avoid a brief flickering of the login view when you're already logged in
    return null;
  }

  if (!user) {
    return <Login />;
  }

  return <>{children}</>;
}
