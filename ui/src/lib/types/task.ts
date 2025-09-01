import type { API } from "./api";
import { Base } from "./general";

export enum TaskStatus {
  RUNNING = "running",
  WAITING = "waiting",
}

export interface Task extends Base {
  name: string;
  status: TaskStatus;
  nextRun: Date;
  recurring: boolean;
  lastStatus?: TaskResult;
  lastRun?: Date;
  lastError?: string;
  interval?: number;
}

export enum TaskResult {
  SUCCESS = "success",
  FAILED = "failed",
  RESOLVED = "resolved",
}

export interface TaskHistory extends Base {
  name: string;
  result: TaskResult;
  runAt: Date;
  error?: string;
  recurring: boolean;
  duration: number;
}

export interface TaskHistoryFilter {
  name?: string;
  result?: TaskResult;
}

export function convertTaskToModel(task: API.Task): Task {
  return {
    id: task.id,
    name: task.name,
    status: task.status as TaskStatus,
    nextRun: new Date(task.next_run),
    recurring: task.recurring,
    lastStatus: task.last_status ? task.last_status as TaskResult : undefined,
    lastRun: task.last_run ? new Date(task.last_run) : undefined,
    lastError: task.last_error,
    interval: task.interval,
  };
}

export function convertTasksToModel(tasks: API.Task[]): Task[] {
  return tasks.map(convertTaskToModel);
}

export function convertTaskHistoryToModel(history: API.TaskHistory[]): TaskHistory[] {
  return history.map(history => ({
    id: history.id,
    name: history.name,
    result: history.result as TaskResult,
    runAt: new Date(history.run_at),
    error: history.error,
    recurring: history.recurring,
    duration: history.duration,
  }));
}
