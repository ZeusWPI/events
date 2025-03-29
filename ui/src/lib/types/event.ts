import type { API } from "./api";
import type { Organizer } from "./organizer";
import type { Year } from "./year";
import { convertOrganizersToJSON, convertOrganizerToModel } from "./organizer";
import { convertYearToJSON, convertYearToModel } from "./year";

export interface Event {
  id: number;
  url: string;
  name: string;
  description: string;
  startTime: Date;
  endTime?: Date;
  location: string;
  year: Year;
  organizers: Organizer[];
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
  };
}

export function convertEventsToModel(events: API.Event[]): Event[] {
  return events.map(convertEventToModel).sort((a, b) => a.startTime.getTime() - b.startTime.getTime());
}

export function convertEventToJSON(event: Event): API.Event {
  return {
    id: event.id,
    url: event.url,
    name: event.name,
    description: event.description,
    start_time: event.startTime.toISOString(),
    end_time: event.endTime?.toISOString() ?? "",
    location: event.location,
    year: convertYearToJSON(event.year),
    organizers: convertOrganizersToJSON(event.organizers),
  };
}
