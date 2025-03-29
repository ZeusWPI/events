import type { API } from "./api";

export interface Organizer {
  id: number;
  role: string;
  name: string;
}

export function convertOrganizerToModel(organizer: API.Organizer): Organizer {
  return {
    id: organizer.id,
    role: organizer.role,
    name: organizer.name,
  };
}

export function convertOrganizersToModel(organizers: API.Organizer[]): Organizer[] {
  return organizers.map(convertOrganizerToModel);
}

export function convertOrganizerToJSON(organizer: Organizer): API.Organizer {
  return {
    id: organizer.id,
    role: organizer.role,
    name: organizer.name,
  };
}

export function convertOrganizersToJSON(organizers: Organizer[]): API.Organizer[] {
  return organizers.map(organizer => convertOrganizerToJSON(organizer));
}
