import { readFileSync } from "fs";
import OpenAI from "openai";
import {
  ChatCompletionMessageParam,
  ChatCompletionTool,
} from "openai/resources";

import { Message, MessageType, UserContext } from "../__generated__/types";
import {
  openaiCreateTimeTracking,
  openaiUpdateTimeTracking,
} from "../functions/statTracking";
import { handleToolCalls } from "./utils";

const openai = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const chatCompletionConfig = {
  model: "gpt-4o-mini",
  temperature: 0.8,
  top_p: 0.9,
  max_completion_tokens: 1024,
};

const instructions = readFileSync("./src/openai/instructions.md", "utf8");

const tools: ChatCompletionTool[] = [
  { type: "function", function: openaiCreateTimeTracking },
  { type: "function", function: openaiUpdateTimeTracking },
];

export const chat = async (props: {
  userContext: UserContext;
  userData: unknown;
  history: ChatCompletionMessageParam[];
  content: string;
}): Promise<Message> => {
  let messages: ChatCompletionMessageParam[] = [
    { role: "system", content: instructions },
    { role: "user", content: JSON.stringify(props.userContext) },
    { role: "user", content: JSON.stringify(props.userData) },
    ...props.history,
    { role: "user", content: props.content },
  ];

  const completion = await openai.chat.completions.create({
    ...chatCompletionConfig,
    messages,
    tools,
  });

  messages.push(completion.choices[0].message);

  if (completion.choices[0].message.tool_calls) {
    const toolCallMessages = await handleToolCalls(
      completion.choices[0].message.tool_calls,
    );
    messages = [...messages, ...toolCallMessages];

    const completion2 = await openai.chat.completions.create({
      ...chatCompletionConfig,
      messages,
      tools,
    });
    messages.push(completion2.choices[0].message);

    return {
      content: completion2.choices[0].message.content || "",
      timestamp: new Date().toISOString(),
      type: MessageType.AiMessage,
    };
  } else {
    return {
      content: completion.choices[0].message.content || "",
      timestamp: new Date().toISOString(),
      type: MessageType.AiMessage,
    };
  }
};
