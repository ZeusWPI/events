import type { API } from "./api";
import type { Organizer } from "./organizer";
import type { Year } from "./year";
import { convertOrganizerToModel } from "./organizer";
import { convertYearToModel } from "./year";
import { Check, convertCheckToModel } from "./check";
import { Announcement, convertAnnouncementToModel } from "./announcement";
import { Base } from "./general";
import { convertPosterToModel, Poster } from "./poster";

export interface Event extends Base {
  url: string;
  name: string;
  description: string;
  startTime: Date;
  endTime?: Date;
  location: string;
  year: Year;
  organizers: Organizer[];
  checks: Check[];
  announcement?: Announcement;
  posters: Poster[];
}

export function convertEventToModel(event: API.Event): Event {
  return {
    id: event.id,
    url: event.url,
    name: event.name,
    description: event.description,
    startTime: new Date(event.start_time),
    endTime: event.end_time ? new Date(event.end_time) : undefined,
    location: event.location,
    year: convertYearToModel(event.year),
    organizers: event.organizers.map(convertOrganizerToModel),
    checks: event.checks ? event.checks.map(convertCheckToModel) : [],
    announcement: event.announcement ? convertAnnouncementToModel(event.announcement) : undefined,
    posters: event.posters.map(convertPosterToModel)
  };
}

export function convertEventsToModel(events: API.Event[]): Event[] {
  return events.map(convertEventToModel).sort((a, b) => a.startTime.getTime() - b.startTime.getTime());
}

