import { API } from "./api";
import { Base } from "./general";

export interface Poster extends Base {
  eventId: number;
  scc: boolean;
}

export function convertPosterToModel(poster: API.Poster): Poster {
  return {
    id: poster.id,
    eventId: poster.event_id,
    scc: poster.scc,
  }
}
