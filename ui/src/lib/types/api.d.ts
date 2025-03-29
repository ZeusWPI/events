import type { JSON } from "./general";

export namespace API {

  interface Base extends JSON {
    id: number;
  }

  export interface Event extends Base {
    url: string;
    name: string;
    description: string;
    start_time: string;
    end_time?: string;
    location: string;
    year: Year;
    organizers: Organizer[];
  }

  export interface Year extends Base {
    start_year: number;
    end_year: number;
  }

  export interface Organizer extends Base {
    role: string;
    name: string;
  }

  export interface Task extends Base {
    name: string;
    status: string;
    next_run: string;
    recurring: boolean;
    last_status?: string;
    last_run?: string;
    last_error?: string;
    interval?: number;
  }

  export interface TaskHistory extends Base {
    name: string;
    result: string;
    run_at: string;
    error?: string;
    recurring: boolean;
  }
}
