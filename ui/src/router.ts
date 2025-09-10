import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
import App from "./App";
import Error404 from "./pages/404";
import { Announcements } from "./pages/announcements/Announcements";
import { AnnouncementsCreate } from "./pages/announcements/AnnouncementsCreate";
import { AnnouncementsOverview } from "./pages/announcements/AnnouncementsOverview";
import { Error } from "./pages/Error";
import Events from "./pages/events/Events";
import { EventsAssign } from "./pages/events/EventsAssign";
import { EventsDetail } from "./pages/events/EventsDetail";
import { EventsOverview } from "./pages/events/EventsOverview";
import Index from "./pages/Index";
import { Mails } from "./pages/mails/Mails";
import { MailsCreate } from "./pages/mails/MailsCreate";
import { PowerPoints } from "./pages/powerpoints/Powerpoints";
import { Tasks } from "./pages/tasks/Tasks";
import { TasksDetail } from "./pages/tasks/TasksDetail";
import { TasksHistory } from "./pages/tasks/TasksHistory";
import { TasksOverview } from "./pages/tasks/TasksOverview";
import { AnnouncementsEdit } from "./pages/announcements/AnnouncementsEdit";
import { MailsEdit } from "./pages/mails/MailsEdit";
import { MailsOverview } from "./pages/mails/MailsOverview";

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

const eventsOverviewRoute = createRoute({
  getParentRoute: () => eventsRoute,
  path: "/",
  component: EventsOverview,
})

const eventsDetailRoute = createRoute({
  getParentRoute: () => eventsRoute,
  path: "/$id",
  component: EventsDetail,
});

const eventsAssignRoute = createRoute({
  getParentRoute: () => eventsRoute,
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

const announcementsOverviewRoute = createRoute({
  getParentRoute: () => announcementsRoute,
  path: "/",
  component: AnnouncementsOverview,
})

const announcementsCreateRoute = createRoute({
  getParentRoute: () => announcementsRoute,
  path: "/create",
  component: AnnouncementsCreate,
})

const announcementsEditRoute = createRoute({
  getParentRoute: () => announcementsRoute,
  path: "/edit/$announcementId",
  component: AnnouncementsEdit,
})

//
// Mails
//

const mailsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/mails",
  component: Mails,
})

const mailsOverviewRoute = createRoute({
  getParentRoute: () => mailsRoute,
  path: "/",
  component: MailsOverview,
})

const mailsCreateRoute = createRoute({
  getParentRoute: () => mailsRoute,
  path: "/create",
  component: MailsCreate,
})

const mailsEditRoute = createRoute({
  getParentRoute: () => mailsRoute,
  path: "/edit/$mailId",
  component: MailsEdit,
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
  path: "/$uid",
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
    eventsOverviewRoute,
    eventsDetailRoute,
    eventsAssignRoute,
  ]),
  announcementsRoute.addChildren([
    announcementsOverviewRoute,
    announcementsCreateRoute,
    announcementsEditRoute,
  ]),
  mailsRoute.addChildren([
    mailsOverviewRoute,
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
