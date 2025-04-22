import { ChatCompletionTool } from "openai/resources/index.mjs";

import { functionGetCalendarEvents, toolGetCalendarEvents } from "./calendar";
import { functionGetCurrentDateTime, toolGetCurrentDateTime } from "./datetime";
import { functionGetMails, toolGetMails } from "./gmail";
import { FunctionName } from "./types";

export * from "./types";

export const tools: ChatCompletionTool[] = [
  toolGetCurrentDateTime,
  toolGetCalendarEvents,
  toolGetMails,
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const callFunction = async (functionName: FunctionName, parameters: any): Promise<any> => {
  switch (functionName) {
    case FunctionName.GetCurrentDateTime:
      return functionGetCurrentDateTime(parameters);
    case FunctionName.GetCalendarEvents:
      return functionGetCalendarEvents(parameters);
    case FunctionName.GetMails:
      return functionGetMails(parameters);
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
