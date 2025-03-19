import type { Organizer } from "../types/types";
import { createContext } from "react";

interface AuthContextType {
  user: Organizer | null;
  isLoading: boolean;
  forbidden: boolean;
  login: () => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
