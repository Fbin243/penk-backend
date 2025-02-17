import { readFileSync } from "fs";
import dotenv from "dotenv";
import readline from "readline";
import OpenAI from "openai";
import chalk from "chalk";
import {
  openaiCreateTimeTracking,
  openaiUpdateTimeTracking,
} from "./functions/statTracking";
import {
  ChatCompletionMessageParam,
  ChatCompletionTool,
} from "openai/resources";
import { User } from "./functions/types";
import { handleToolCalls } from "./utils";

dotenv.config({
  path: ".env.penk",
});

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

const openai = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const chatCompletionConfig = {
  model: "gpt-4o-mini",
  temperature: 0.8,
  top_p: 0.9,
  max_completion_tokens: 1024,
};

const instructions = readFileSync("./instructions.md", "utf8");

const testUser: User = {
  id: "test-user-123",
  name: "Test User",
  currentCharacterId: "test-char-1",
  characters: [
    {
      id: "test-char-1",
      name: "Test Char 1",
      categories: [
        {
          id: "ct1",
          name: "books",
        },
        {
          id: "ct2",
          name: "music",
        },
      ],
    },
    {
      id: "test-char-2",
      name: "Test Char 2",
      categories: [
        {
          id: "ct3",
          name: "coding",
        },
      ],
    },
  ],
};
console.log(chalk.redBright("[Inject user context]"));
console.dir(testUser, { depth: null, colors: true });
console.log();

let messages: ChatCompletionMessageParam[] = [
  { role: "system", content: instructions },
  { role: "user", content: JSON.stringify(testUser) },
];

const tools: ChatCompletionTool[] = [
  { type: "function", function: openaiCreateTimeTracking },
  { type: "function", function: openaiUpdateTimeTracking },
];

const startChat = () => {
  rl.question("You: ", async (input) => {
    if (input.toLowerCase() === "exit") {
      console.log("Goodbye!");
      rl.close();
    } else {
      messages.push({ role: "user", content: input });

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
      } else {
        messages.push(completion.choices[0].message);
      }

      console.log(chalk.cyan(`Bot: ${messages[messages.length - 1].content}`));
      startChat();
    }
  });
};

console.log(chalk.yellowBright('Welcome to the chat! Type "exit" to quit.'));
startChat();
