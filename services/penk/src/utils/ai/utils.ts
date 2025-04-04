import { readFileSync } from "fs";
import OpenAI from "openai";
import {
  ChatCompletionAssistantMessageParam,
  ChatCompletionChunk,
  ChatCompletionMessageParam,
  ChatCompletionUserMessageParam,
  CompletionUsage,
} from "openai/resources/index.mjs";

import { Message, MessageType } from "../types";
import { tools } from "./functions";

const baseInstruction = readFileSync("resources/instructions/base.md", "utf8");

export const setupInitialMessages = (
  userData: string,
  history: Message[],
): ChatCompletionMessageParam[] => {
  return [
    {
      role: "system",
      content: baseInstruction,
    },
    {
      role: "user",
      content: `Current date and time: ${new Date().toISOString()}`,
    },
    {
      role: "user",
      content: `My data: ${userData}`,
    },
    ...history.map((message) => {
      if (message.type === MessageType.UserMessage) {
        return {
          role: "user",
          content: message.content,
        } satisfies ChatCompletionUserMessageParam;
      }
      return {
        role: "assistant",
        content: message.content,
      } satisfies ChatCompletionAssistantMessageParam;
    }),
  ];
};

export const streamAssistantResponse = async (props: {
  client: OpenAI;
  messages: ChatCompletionMessageParam[];
  aiTimestamp: string;
  onChunk: (chunk: string, timestamp: string) => void;
}): Promise<{
  content: string;
  toolCalls: Record<number, ChatCompletionChunk.Choice.Delta.ToolCall>;
  usage?: CompletionUsage;
}> => {
  const { client, messages, aiTimestamp, onChunk } = props;

  let content = "";
  let usage: CompletionUsage | undefined;
  const finalToolCalls: Record<number, ChatCompletionChunk.Choice.Delta.ToolCall> = {};

  const stream = await client.chat.completions.create({
    model: "gpt-4o-mini",
    messages,
    stream: true,
    stream_options: { include_usage: true },
    max_tokens: 2048,
    tools,
    tool_choice: "auto",
  });

  for await (const chunk of stream) {
    const toolCalls = chunk.choices[0]?.delta?.tool_calls || [];
    for (const toolCall of toolCalls) {
      const { index } = toolCall;

      if (!finalToolCalls[index]) {
        finalToolCalls[index] = toolCall;
      }

      if (finalToolCalls[index].function && toolCall.function?.arguments) {
        finalToolCalls[index].function.arguments += toolCall.function.arguments;
      }
    }

    const partial = chunk.choices[0]?.delta?.content || "";
    if (partial) {
      content += partial;
      onChunk(partial, aiTimestamp);
    }

    if (chunk.usage) {
      usage = chunk.usage;
    }
  }

  return { content, toolCalls: finalToolCalls, usage };
};
