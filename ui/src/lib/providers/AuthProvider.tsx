import type { ReactNode } from "react";
import type { Organizer } from "../types/types";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { useUser, useUserLogout } from "../api/user";
import { AuthContext } from "../context/authContext";

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<Organizer | null>(null);

  const { data, isLoading } = useUser();
  const { mutate: logoutMutation } = useUserLogout();

  const url = import.meta.env.VITE_BACKEND_URL as string;

  useEffect(() => {
    if (data) {
      setUser(data);
    }
  }, [data]);

  const login = useCallback(() => {
    window.location.href = `${url}/auth/login/zauth`;
  }, [url]);

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: () => setUser(null),
      onError: (err) => {
        toast.error("Logout failed");
        console.error(err);
      },
    });
  }, [logoutMutation]);

  const value = useMemo(() => ({ user, isLoading, login, logout }), [user, isLoading, login, logout]);

  return <AuthContext value={value}>{children}</AuthContext>;
}
