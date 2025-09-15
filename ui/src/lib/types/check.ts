import { SquareCheckIcon, SquareIcon, SquareXIcon, TriangleAlertIcon } from "lucide-react";
import { createElement } from "react";
import { API } from "./api";
import { Base } from "./general";

export enum CheckStatus {
  Done = "done",
  DoneLate = "done_late",
  Todo = "todo",
  TodoLate = "todo_late",
  Warning = "warning",
}

export const checkStatusToIcon: Record<CheckStatus, React.ReactNode> = {
  [CheckStatus.Done]: createElement(SquareCheckIcon, { className: 'text-green-500' }),
  [CheckStatus.DoneLate]: createElement(SquareCheckIcon, { className: 'text-orange-500' }),
  [CheckStatus.Todo]: createElement(SquareIcon),
  [CheckStatus.TodoLate]: createElement(SquareXIcon, { className: 'text-red-500' }),
  [CheckStatus.Warning]: createElement(TriangleAlertIcon, { className: 'text-orange-500' }),
}

export const checkStatusToText: Record<CheckStatus, string> = {
  [CheckStatus.Done]: "Done",
  [CheckStatus.DoneLate]: "Done after deadline",
  [CheckStatus.Todo]: "Todo",
  [CheckStatus.TodoLate]: "Deadline passed",
  [CheckStatus.Warning]: "Warning",
}

export enum CheckType {
  Manual = "manual",
  Automatic = "automatic",
}

export interface Check extends Base {
  eventId: number;
  status: CheckStatus;
  message?: string;
  description: string;
  deadline?: number;
  type: CheckType;
  creator_id?: number;
}

export function convertCheckToModel(check: API.Check): Check {
  return {
    id: check.id,
    eventId: check.event_id,
    status: check.status as CheckStatus,
    message: check.message,
    description: check.description,
    deadline: check.deadline,
    type: check.type as CheckType,
    creator_id: check.creator_id,
  };
}

