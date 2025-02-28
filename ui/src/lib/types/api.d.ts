export namespace API {

  interface Base {
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
}
