import { ChatCompletionTool } from "openai/resources/index.mjs";

import { Tool } from "../../types";
import { functionGetCalendarEvents, toolGetCalendarEvents } from "./calendar";
import { functionGetMails, toolGetMails } from "./gmail";

export const tools: ChatCompletionTool[] = [toolGetCalendarEvents, toolGetMails];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const callFunction = async (functionName: Tool, parameters: any): Promise<any> => {
  switch (functionName) {
    case Tool.GetCalendarEvents:
      return functionGetCalendarEvents(parameters);
    case Tool.GetEmails:
      return functionGetMails(parameters);
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
