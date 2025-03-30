import chalk from "chalk";
import { readFileSync } from "fs";
import OpenAI from "openai";
import {
  ChatCompletionAssistantMessageParam,
  ChatCompletionMessageParam,
  ChatCompletionUserMessageParam,
  CompletionUsage,
} from "openai/resources/index.mjs";

import { Message, MessageType } from "../../utils/types";

const client = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const instructions = readFileSync("instructions.md", "utf8");

const analyzeUsage = (usage: CompletionUsage) => {
  /*
    Price
    - Input: $0.150 / 1M tokens
    - Cached input: $0.075 / 1M tokens
    - Output: $0.600 / 1M tokens
  */

  const cachedTokens = usage.prompt_tokens_details?.cached_tokens || 0;
  const inputTokens = usage.prompt_tokens - cachedTokens;
  const inputCost = (inputTokens * 0.15) / 1000000;
  const cachedInputCost = (cachedTokens * 0.075) / 1000000;
  const outputCost = (usage.completion_tokens * 0.6) / 1000000;

  return {
    cachedTokens,
    inputTokens,
    completionTokens: usage.completion_tokens,
    totalCost: inputCost + cachedInputCost + outputCost,
  };
};

export const chat = async (props: { userData: string; messages: Message[] }): Promise<Message> => {
  const openAiMessages: ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: instructions,
    },
    {
      role: "user",
      content: props.userData,
    },
    ...props.messages.map((message) => {
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

  const completion = await client.chat.completions.create({
    model: "gpt-4o-mini",
    messages: openAiMessages,
  });

  // Log the usage
  if (completion.usage) {
    console.log(chalk.green("Usage:"));
    console.log(analyzeUsage(completion.usage));
    console.log();
  }

  return {
    type: MessageType.AiMessage,
    content: completion.choices[0].message.content || "",
    timestamp: new Date().toISOString(),
  };
};
