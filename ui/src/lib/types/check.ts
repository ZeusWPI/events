import { CheckIcon, XIcon, CircleAlertIcon } from "lucide-react";
import { createElement } from "react";
import { randomNumber } from "../utils/utils";
import { API } from "./api";
import { Base } from "./general";

export enum CheckStatus {
  Finished = "finished",
  Unfinished = "unfinished",
  Warning = "warning",
}

export const statusToIcon: Record<CheckStatus, React.ReactNode> = {
  [CheckStatus.Finished]: createElement(CheckIcon, { className: 'text-green-500' }),
  [CheckStatus.Unfinished]: createElement(XIcon, { className: 'text-red-500' }),
  [CheckStatus.Warning]: createElement(CircleAlertIcon, { className: 'text-orange-500' }),
}

export enum CheckSource {
  Automatic = "automatic",
  Manual = "manual",
}

export interface Check extends Base {
  eventId: number;
  description: string;
  warning?: string;
  status: CheckStatus;
  error?: string;
  source: CheckSource;
}

export function convertCheckToModel(check: API.Check): Check {
  return {
    id: check.source as CheckSource === CheckSource.Manual ? check.id : randomNumber(),
    eventId: check.event_id,
    description: check.description,
    warning: check.warning,
    status: check.status as CheckStatus,
    error: check.error,
    source: check.source as CheckSource,
  };
}

