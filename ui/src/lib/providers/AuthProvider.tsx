import type { ReactNode } from "react";
import type { Organizer } from "../types/types";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { useUser, useUserLogout } from "../api/user";
import { AuthContext } from "../context/authContext";
import { isResponseNot200Error } from "../utils/query";

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<Organizer | null>(null);
  const [forbidden, setForbidden] = useState(false);

  const { data, isLoading, error } = useUser();
  const { mutate: logoutMutation } = useUserLogout();

  useEffect(() => {
    if (data) {
      setUser(data);
      setForbidden(false);
    }
  }, [data]);

  useEffect(() => {
    if (error && isResponseNot200Error(error)) {
      if (error.response.status === 403) {
        setForbidden(true);
        return;
      }
    }

    setForbidden(false);
  }, [error]);

  const url = import.meta.env.VITE_BACKEND_URL as string;

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

  const value = useMemo(() => ({ user, isLoading, forbidden, login, logout }), [user, isLoading, forbidden, login, logout]);

  return <AuthContext value={value}>{children}</AuthContext>;
}
