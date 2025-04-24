import { ChatCompletionTool } from "openai/resources/index.mjs";

import { Tool } from "../../types";
import { functionGetCalendarEvents, toolGetCalendarEvents } from "./calendar";
import { functionGetMails, toolGetMails } from "./gmail";
import { functionGetHabits, toolGetHabits } from "./habits";
import { functionGetMetrics, toolGetMetrics } from "./metrics";
import {
  functionCreateTask,
  functionCreateTaskSession,
  functionDeleteTask,
  functionDeleteTaskSession,
  functionGetTasks,
  functionPlanDay,
  functionUpdateTask,
  functionUpdateTaskSession,
  toolCreateTask,
  toolCreateTaskSession,
  toolDeleteTask,
  toolDeleteTaskSession,
  toolGetTasks,
  toolPlanDay,
  toolUpdateTask,
  toolUpdateTaskSession,
} from "./tasks";

export const tools: ChatCompletionTool[] = [
  toolGetCalendarEvents,
  toolGetMails,
  toolGetTasks,
  toolCreateTask,
  toolUpdateTask,
  toolDeleteTask,
  toolCreateTaskSession,
  toolUpdateTaskSession,
  toolDeleteTaskSession,
  toolPlanDay,
  toolGetMetrics,
  toolGetHabits,
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const callFunction = async (functionName: Tool, parameters: any): Promise<any> => {
  switch (functionName) {
    case Tool.GetCalendarEvents:
      return functionGetCalendarEvents(parameters);
    case Tool.GetEmails:
      return functionGetMails(parameters);
    case Tool.GetTasks:
      return functionGetTasks(parameters);
    case Tool.CreateTask:
      return functionCreateTask(parameters);
    case Tool.UpdateTask:
      return functionUpdateTask(parameters);
    case Tool.DeleteTask:
      return functionDeleteTask(parameters);
    case Tool.CreateTaskSession:
      return functionCreateTaskSession(parameters);
    case Tool.UpdateTaskSession:
      return functionUpdateTaskSession(parameters);
    case Tool.DeleteTaskSession:
      return functionDeleteTaskSession(parameters);
    case Tool.PlanDay:
      return functionPlanDay(parameters);
    case Tool.GetMetrics:
      return functionGetMetrics(parameters);
    case Tool.GetHabits:
      return functionGetHabits(parameters);
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
