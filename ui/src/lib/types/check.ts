import { randomNumber } from "../utils/utils";
import { API } from "./api";

export enum CheckSource {
  Automatic = "automatic",
  Manual = "manual",
}

export interface Check {
  id: number;
  eventId: number;
  description: string;
  done: boolean;
  error?: string;
  source: CheckSource;
}

export function convertCheckToModel(check: API.Check): Check {
  console.log(check)
  return {
    id: check.source as CheckSource === CheckSource.Manual ? check.id : randomNumber(),
    eventId: check.event_id,
    description: check.description,
    done: check.done,
    error: check.error,
    source: check.source as CheckSource,
  };
}

