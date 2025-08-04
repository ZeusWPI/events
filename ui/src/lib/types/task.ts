import type { API } from "./api";
import { Base } from "./general";

export enum TaskStatus {
  running = "running",
  waiting = "waiting",
}

export interface Task extends Base {
  name: string;
  status: TaskStatus;
  nextRun: Date;
  recurring: boolean;
  lastStatus?: TaskHistoryStatus;
  lastRun?: Date;
  lastError?: string;
  interval?: number;
}

export enum TaskHistoryStatus {
  success = "success",
  failed = "failed",
}

export interface TaskHistory extends Base {
  name: string;
  result: TaskHistoryStatus;
  runAt: Date;
  error?: string;
  recurring: boolean;
  duration: number;
}

export interface TaskHistoryFilter {
  name?: string;
  onlyErrored?: boolean;
  recurring?: boolean;
}

export function convertTaskStatusToModel(status: string): TaskStatus {
  if (Object.values(TaskStatus).includes(status as TaskStatus)) {
    return status as TaskStatus;
  }

  // Can only happen if the backend and frontend statuses are out of sync
  return TaskStatus.waiting;
}

export function convertTaskToModel(task: API.Task): Task {
  return {
    id: task.id,
    name: task.name,
    status: convertTaskStatusToModel(task.status),
    nextRun: new Date(task.next_run),
    recurring: task.recurring,
    lastStatus: task.last_status ? convertTaskHistoryStatusToModel(task.last_status) : undefined,
    lastRun: task.last_run ? new Date(task.last_run) : undefined,
    lastError: task.last_error,
    interval: task.interval,
  };
}

export function convertTasksToModel(tasks: API.Task[]): Task[] {
  return tasks.map(convertTaskToModel);
}

export function convertTaskHistoryStatusToModel(status: string): TaskHistoryStatus {
  if (Object.values(TaskHistoryStatus).includes(status as TaskHistoryStatus)) {
    return status as TaskHistoryStatus;
  }

  // Can only happend if the backend and frontend statuses are out of sync
  return TaskHistoryStatus.failed;
}

export function convertTaskHistoryToModel(history: API.TaskHistory[]): TaskHistory[] {
  return history.map(history => ({
    id: history.id,
    name: history.name,
    result: convertTaskHistoryStatusToModel(history.result),
    runAt: new Date(history.run_at),
    error: history.error,
    recurring: history.recurring,
    duration: history.duration,
  }));
}
