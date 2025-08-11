import type { Year } from "../types/year";

import { createContext } from "react";

interface YearContextType {
  year: Year;
  setYear: (year: Year) => void;
  isLoading: boolean;
  locked: boolean;
  setLocked: (locked: boolean) => void;
}

const initialState: YearContextType = {
  year: { id: 0, start: 0, end: 0, formatted: "" },
  setYear: () => null,
  isLoading: true,
  locked: false,
  setLocked: () => null,
}

export const YearContext = createContext<YearContextType>(initialState)
