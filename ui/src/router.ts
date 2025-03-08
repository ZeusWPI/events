import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
import App from "./App";
import Error404 from "./pages/404";
import Events from "./pages/events/Events";
import { EventsAssign } from "./pages/events/EventsAssign";
import { EventsDetail } from "./pages/events/EventsDetail";
import { EventsYear } from "./pages/events/EventsYear";
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

const eventsYearRoute = createRoute({
  getParentRoute: () => eventsRoute,
  path: "$year",
  component: EventsYear,
});

const eventsYearDetailRoute = createRoute({
  getParentRoute: () => eventsYearRoute,
  path: "/$id",
  component: EventsDetail,
});

const eventsYearAssignRoute = createRoute({
  getParentRoute: () => eventsYearRoute,
  path: "/assign",
  component: EventsAssign,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  eventsRoute.addChildren([
    eventsYearRoute.addChildren([eventsYearDetailRoute, eventsYearAssignRoute]),
  ]),
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
