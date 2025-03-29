import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
import App from "./App";
import Error404 from "./pages/404";
import { Error } from "./pages/Error";
import Events from "./pages/events/Events";
import { EventsAssign } from "./pages/events/EventsAssign";
import { EventsDetail } from "./pages/events/EventsDetail";
import { EventsYear } from "./pages/events/EventsYear";
import Index from "./pages/Index";
import { Tasks } from "./pages/tasks/Tasks";
import { TasksDetail } from "./pages/tasks/TasksDetail";
import { TasksHistory } from "./pages/tasks/TasksHistory";
import { TasksOverview } from "./pages/tasks/TasksOverview";

const rootRoute = createRootRoute({
  component: App,
  notFoundComponent: Error404,
  errorComponent: Error,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: Index,
});

//
// Events
//

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

//
// Tasks
//

const tasksRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/tasks",
  component: Tasks,
});

const tasksOverviewRoute = createRoute({
  getParentRoute: () => tasksRoute,
  path: "/",
  component: TasksOverview,
});

const tasksDetailRoute = createRoute({
  getParentRoute: () => tasksRoute,
  path: "/$id",
  component: TasksDetail,
});

const tasksHistory = createRoute({
  getParentRoute: () => tasksRoute,
  path: "/history",
  component: TasksHistory,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  eventsRoute.addChildren([
    eventsYearRoute.addChildren([eventsYearDetailRoute, eventsYearAssignRoute]),
  ]),
  tasksRoute.addChildren([tasksOverviewRoute, tasksDetailRoute, tasksHistory]),
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
