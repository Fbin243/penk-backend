import { readFileSync } from "fs";
import OpenAI from "openai";
import {
  ChatCompletionAssistantMessageParam,
  ChatCompletionMessageParam,
  ChatCompletionUserMessageParam,
} from "openai/resources/index.mjs";

import { Message, MessageType } from "../../utils/types";

const client = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const instructions = readFileSync("instructions.md", "utf8");

export const chat = async (messages: Message[]) => {
  const openAiMessages: ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: instructions,
    },
    ...messages.map((message) => {
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

  return completion.choices[0].message.content;
};
