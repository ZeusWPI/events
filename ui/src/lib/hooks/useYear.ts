import { use, useEffect } from "react";
import { YearContext } from "../contexts/yearContext";

export function useYear() {
  const context = use(YearContext);
  if (!context) {
    throw new Error("useYear must be used within a YearProvider")
  }
  return context;
}

export function useYearLock() {
  const { setLocked } = useYear();

  useEffect(() => {
    setLocked(true);

    return () => {
      setLocked(false);
    }
  }, [setLocked]);
}
