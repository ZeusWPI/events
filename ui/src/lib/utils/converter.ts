import type { API } from "../types/api";
import type { Base, Event, Organizer, Year } from "../types/types";

// API to internal model

function convertBaseToModel(item: API.Base): Base {
  return {
    id: item.id,
  };
}

export function convertEventToModel(event: API.Event): Event {
  return {
    ...convertBaseToModel(event),
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

export function convertYearToModel(year: API.Year): Year {
  const startFormatted = (year.start_year % 100).toString().padStart(2, "0");
  const endFormatted = (year.end_year % 100).toString().padStart(2, "0");

  return {
    ...convertBaseToModel(year),
    startYear: year.start_year,
    endYear: year.end_year,
    formatted: `${startFormatted}-${endFormatted}`,
  };
}

export function convertYearsToModel(years: API.Year[]): Year[] {
  return years.map(convertYearToModel);
}

export function convertOrganizerToModel(organizer: API.Organizer): Organizer {
  return {
    ...convertBaseToModel(organizer),
    role: organizer.role,
    name: organizer.name,
  };
}

export function convertOrganizersToModel(organizers: API.Organizer[]): Organizer[] {
  return organizers.map(convertOrganizerToModel);
}

// Internal model to API

export function convertBaseToJSON(item: Base): API.Base {
  return {
    id: item.id,
  };
}

export function convertEventToJSON(event: Event): API.Event {
  return {
    ...convertBaseToJSON(event),
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

export function convertYearToJSON(year: Year): API.Year {
  return {
    ...convertBaseToJSON(year),
    start_year: year.startYear,
    end_year: year.endYear,
  };
}

export function convertOrganizerToJSON(organizer: Organizer): API.Organizer {
  return {
    ...convertBaseToJSON(organizer),
    role: organizer.role,
    name: organizer.name,
  };
}

export function convertOrganizersToJSON(organizers: Organizer[]): API.Organizer[] {
  return organizers.map(organizer => convertOrganizerToJSON(organizer));
}
