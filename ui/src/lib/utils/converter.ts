import type { API } from "../types/api";
import type { Base, Event, Organizer, Year } from "../types/types";

// API Converters

function convertBase(items: API.Base): Base {
  return {
    id: items.id,
  };
}

export function convertEvent(event: API.Event): Event {
  return {
    ...convertBase(event),
    url: event.url,
    name: event.name,
    description: event.description,
    startTime: new Date(event.start_time),
    endTime: event.end_time ? new Date(event.end_time) : undefined,
    location: event.location,
    year: convertYear(event.year),
    organizers: event.organizers.map(convertOrganizer),
  };
}

export function convertEvents(events: API.Event[]): Event[] {
  return events.map(convertEvent).sort((a, b) => a.startTime.getTime() - b.startTime.getTime());
}

export function convertYear(year: API.Year): Year {
  return {
    ...convertBase(year),
    startYear: year.start_year,
    endYear: year.end_year,

  };
}

export function convertYears(years: API.Year[]): Year[] {
  return years.map(convertYear);
}

export function convertOrganizer(organizer: API.Organizer): Organizer {
  return {
    ...convertBase(organizer),
    role: organizer.role,
    name: organizer.name,
  };
}

export function convertOrganizers(organizers: API.Organizer[]): Organizer[] {
  return organizers.map(convertOrganizer);
}

// Other

export function yearToString({ startYear, endYear }: Year): string {
  const startFormatted = (startYear % 100).toString().padStart(2, "0");
  const endFormatted = (endYear % 100).toString().padStart(2, "0");

  return `${startFormatted}-${endFormatted}`;
}
