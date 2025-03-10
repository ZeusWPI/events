import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "./components/ui/sonner.tsx";
import { BreadcrumbProvider } from "./lib/providers/BreadcrumbProvider.tsx";
import { ThemeProvider } from "./lib/providers/ThemeProvider.tsx";
import { router } from "./router.ts";
import "./index.css";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <BreadcrumbProvider>
          <Toaster richColors toastOptions={{}} />
          <RouterProvider router={router} />
        </BreadcrumbProvider>
      </ThemeProvider>
    </QueryClientProvider>
  </StrictMode>,
);
