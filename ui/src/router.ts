import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
import App from "./App";
import Error404 from "./pages/404";
import Events from "./pages/Events";
import Index from "./pages/Index";

const rootRoute = createRootRoute({
  component: App,
  notFoundComponent: Error404,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: Index,
});

const eventsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/events",
  component: Events,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  eventsRoute,
]);

export const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  defaultStaleTime: 5000,
  scrollRestoration: true,
});

// Ensure type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
