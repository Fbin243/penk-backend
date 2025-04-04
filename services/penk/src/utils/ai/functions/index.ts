import { ChatCompletionTool } from "openai/resources/index.mjs";

import { functionGetCalendarEvents, toolGetCalendarEvents } from "./getCalendarEvents";
import { FunctionName } from "./types";

export * from "./types";

export const tools: ChatCompletionTool[] = [toolGetCalendarEvents];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const callFunction = async (functionName: FunctionName, parameters: any): Promise<any> => {
  switch (functionName) {
    case FunctionName.GetCalendarEvents:
      return functionGetCalendarEvents(parameters);
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
