import type { API } from "./api";
import { Base } from "./general";

export enum TaskStatus {
  Waiting = "waiting",
  Running = "running",
}

export enum TaskResult {
  Succes = "success",
  Failed = "failed",
  Resolved = "resolved",
}

export enum TaskType {
  Recurring = "recurring",
  Once = "once",
}

export interface Task {
  uid: string;
  name: string;
  status: TaskStatus;
  nextRun: Date;
  type: TaskType;
  lastStatus?: TaskResult;
  lastRun?: Date;
  lastError?: string;
  interval?: number;
}

export interface TaskHistory extends Base {
  name: string;
  result: TaskResult;
  runAt: Date;
  error?: string;
  type: TaskType;
  duration: number;
}

export interface TaskHistoryFilter {
  uid?: string;
  result?: TaskResult;
}

export function convertTaskToModel(task: API.Task): Task {
  return {
    uid: task.uid,
    name: task.name,
    status: task.status as TaskStatus,
    nextRun: new Date(task.next_run),
    type: task.type as TaskType,
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
    type: history.type as TaskType,
    duration: history.duration,
  }));
}
