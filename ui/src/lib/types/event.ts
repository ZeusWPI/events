import { Announcement, convertAnnouncementToModel } from "./announcement";
import type { API } from "./api";
import { Check, convertCheckToModel } from "./check";
import { Base } from "./general";
import { convertMailToModel, Mail } from "./mail";
import type { Organizer } from "./organizer";
import { convertOrganizerToModel } from "./organizer";
import { convertPosterToModel, Poster } from "./poster";
import type { Year } from "./year";
import { convertYearToModel } from "./year";

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
  announcements: Announcement[];
  mails: Mail[];
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
    organizers: event.organizers?.map(convertOrganizerToModel) ?? [],
    checks: event.checks ? event.checks.map(convertCheckToModel).sort((a, b) => a.description.localeCompare(b.description)) : [],
    announcements: event.announcements?.map(convertAnnouncementToModel).sort((a, b) => a.sendTime.getTime() - b.sendTime.getTime()) ?? [],
    mails: event.mails?.map(convertMailToModel).sort((a, b) => a.sendTime.getTime() - b.sendTime.getTime()) ?? [],
    posters: event.posters?.map(convertPosterToModel) ?? [],
  };
}

export function convertEventsToModel(events: API.Event[]): Event[] {
  return events.map(convertEventToModel).sort((a, b) => a.startTime.getTime() - b.startTime.getTime());
}

