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
    checks?: Check[];
    announcements?: Announcement[];
    mails?: Mail[];
    posters: Poster[];
  }

  export interface Check extends Base {
    event_id: number;
    description: string;
    status: string;
    error: string;
    source: string;
  }

  export interface Announcement extends Base {
    year_id: number;
    event_ids: number[];
    author_id: number;
    content: string;
    send_time: string;
    send: boolean;
    error?: string;
  }

  export interface Year extends Base {
    start: number;
    end: number;
  }

  export interface Organizer extends Base {
    role: string;
    name: string;
    zauth_id: number;
  }

  export interface Poster extends Base {
    event_id: number;
    scc: boolean;
  }

  export interface Mail extends Base {
    year_id: number;
    event_ids: number[];
    author_id: number;
    title: string;
    content: string;
    send_time: string;
    send: boolean;
    error?: string;
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
    duration: number;
  }
}
