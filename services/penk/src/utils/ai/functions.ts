import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";

export enum FunctionName {
  GetDateTime = "get_date_time",
}

const getDateTimeDefinition: FunctionDefinition = {
  name: FunctionName.GetDateTime,
  description: "Get the current date and time",
  parameters: {
    type: "object",
    properties: {},
    additionalProperties: false,
  },
  strict: true,
};

const getDateTime = async () => {
  const dateTime = new Date();
  return dateTime.toLocaleString();
};

export const tools: ChatCompletionTool[] = [
  {
    type: "function",
    function: getDateTimeDefinition,
  },
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unused-vars
export const callFunction = async (functionName: FunctionName, parameters: any): Promise<any> => {
  switch (functionName) {
    case FunctionName.GetDateTime:
      return getDateTime();
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
