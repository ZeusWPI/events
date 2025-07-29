import { randomNumber } from "../utils/utils";
import { API } from "./api";
import { Base } from "./general";

export enum CheckStatus {
  Finished = "finished",
  Unfinished = "unfinished",
  Warning = "warning",
}

export enum CheckSource {
  Automatic = "automatic",
  Manual = "manual",
}

export interface Check extends Base {
  eventId: number;
  description: string;
  status: CheckStatus;
  error?: string;
  source: CheckSource;
}

export function convertCheckToModel(check: API.Check): Check {
  return {
    id: check.source as CheckSource === CheckSource.Manual ? check.id : randomNumber(),
    eventId: check.event_id,
    description: check.description,
    status: check.status as CheckStatus,
    error: check.error,
    source: check.source as CheckSource,
  };
}

