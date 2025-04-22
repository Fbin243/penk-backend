import chalk from "chalk";
import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";

import { SharedDescription } from "./shared";
import { FunctionName } from "./types";

const getCurrentDateTimeDefinition: FunctionDefinition = {
  name: FunctionName.GetCurrentDateTime,
  description:
    "Gets the current date and time formatted according to the user's timezone and locale preferences. Use this tool when the user asks about the current time, date, or needs time-related information in their local format. The function requires both timezone (e.g., 'America/New_York', 'Europe/London') and locale (e.g., 'en-US', 'fr-FR') parameters to correctly format the date and time according to the user's regional settings.",
  parameters: {
    type: "object",
    properties: {
      timezone: {
        type: "string",
        description: SharedDescription.timezone,
      },
      locale: {
        type: "string",
        description: SharedDescription.locale,
      },
    },
    required: ["timezone", "locale"],
    additionalProperties: false,
  },
  strict: true,
};

export const functionGetCurrentDateTime = async (props: { timezone: string; locale: string }) => {
  console.log(chalk.cyan(`[Tool: ${FunctionName.GetCurrentDateTime}]`));
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Date().toLocaleString(props.locale, { timeZone: props.timezone });
};

export const toolGetCurrentDateTime: ChatCompletionTool = {
  type: "function",
  function: getCurrentDateTimeDefinition,
};
