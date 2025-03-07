interface Base {
  id: number;
}

export interface Event extends Base {
  url: string;
  name: string;
  description: string;
  startTime: Date;
  endTime?: Date;
  location: string;
  year: Year;
  organizers: Organizer[];
}

export interface Year extends Base {
  startYear: number;
  endYear: number;
  formatted: string;
}

export interface Organizer extends Base {
  role: string;
  name: string;
}
