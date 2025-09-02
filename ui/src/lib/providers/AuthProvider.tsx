import type { ReactNode } from "react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { isResponseNot200Error } from "../api/query";
import { useUser, useUserLogin, useUserLogout } from "../api/user";
import { AuthContext } from "../contexts/authContext";
import type { Organizer } from "../types/organizer";

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

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: () => toast.success("Logged out"),
      onError: (err) => {
        toast.error("Logout failed");
        console.error(err);
      },
      onSettled: () => setUser(null),
    });
  }, [logoutMutation]);

  const value = useMemo(() => ({ user, isLoading, forbidden, login: useUserLogin, logout }), [user, isLoading, forbidden, logout]);

  return <AuthContext value={value}>{children}</AuthContext>;
}
