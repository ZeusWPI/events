import { use } from "react";
import { ThemeProviderContext } from "../contexts/themeContext";

export function useTheme() {
  const context = use(ThemeProviderContext);

  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }

  let actualTheme = context.theme
  if (context.theme == "system") {
    actualTheme = window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  }

  return { ...context, actualTheme };
}
