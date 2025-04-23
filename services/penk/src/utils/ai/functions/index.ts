import { ChatCompletionTool } from "openai/resources/index.mjs";

import { Tool } from "../../types";
import { functionGetCalendarEvents, toolGetCalendarEvents } from "./calendar";
import { functionGetMails, toolGetMails } from "./gmail";
import {
  functionCreateTask,
  functionCreateTaskSession,
  functionUpdateTask,
  functionUpdateTaskSession,
  toolCreateTask,
  toolCreateTaskSession,
  toolUpdateTask,
  toolUpdateTaskSession,
} from "./tasks";

export const tools: ChatCompletionTool[] = [
  toolGetCalendarEvents,
  toolGetMails,
  toolCreateTask,
  toolUpdateTask,
  toolCreateTaskSession,
  toolUpdateTaskSession,
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const callFunction = async (functionName: Tool, parameters: any): Promise<any> => {
  switch (functionName) {
    case Tool.GetCalendarEvents:
      return functionGetCalendarEvents(parameters);
    case Tool.GetEmails:
      return functionGetMails(parameters);
    case Tool.CreateTask:
      return functionCreateTask(parameters);
    case Tool.UpdateTask:
      return functionUpdateTask(parameters);
    case Tool.CreateTaskSession:
      return functionCreateTaskSession(parameters);
    case Tool.UpdateTaskSession:
      return functionUpdateTaskSession(parameters);
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
