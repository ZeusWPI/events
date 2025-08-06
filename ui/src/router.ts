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
import { Announcements } from "./pages/announcements/Announcements";
import { AnnouncementsYear } from "./pages/announcements/AnnouncementsYear";
import { AnnouncementsCreate } from "./pages/announcements/AnnouncementsCreate";
import { Mails } from "./pages/mails/Mails";
import { MailsCreate } from "./pages/mails/MailsCreate";
import { PowerPoints } from "./pages/powerpoints/Powerpoints";

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
// Announcements
//

const announcementsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/announcements",
  component: Announcements,
})

const announcementsYearRoute = createRoute({
  getParentRoute: () => announcementsRoute,
  path: "$year",
  component: AnnouncementsYear
})

const announcementsCreateRoute = createRoute({
  getParentRoute: () => announcementsYearRoute,
  path: "/$event",
  component: AnnouncementsCreate,
})

//
// Mails
//

const mailsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/mails",
  component: Mails,
})

const mailsCreateRoute = createRoute({
  getParentRoute: () => mailsRoute,
  path: "/create",
  component: MailsCreate,
})

const mailsEditRoute = createRoute({
  getParentRoute: () => mailsRoute,
  path: "/edit/$mail",
  component: MailsCreate
})

//
// Powerpoint
//

const powerpointsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/powerpoints",
  component: PowerPoints,
})

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
  announcementsRoute.addChildren([
    announcementsYearRoute.addChildren([announcementsCreateRoute])
  ]),
  mailsRoute.addChildren([
    mailsCreateRoute,
    mailsEditRoute,
  ]),
  powerpointsRoute,
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
