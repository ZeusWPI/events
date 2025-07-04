import type { API } from "./api";
import { Base } from "./general";

export interface Organizer extends Base {
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
