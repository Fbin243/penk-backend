import "./bootstrap";

import chalk from "chalk";
import { readFileSync } from "fs";

import { chat } from "./services/ai/utils";
import { MessageType } from "./utils/types";

const instructions = readFileSync("instructions.md", "utf8");

const main = async () => {
  const response = await chat([
    {
      type: MessageType.AiMessage,
      content: instructions,
      timestamp: new Date().toISOString(),
    },
    {
      type: MessageType.UserMessage,
      content: "Hello, how are you?",
      timestamp: new Date().toISOString(),
    },
  ]);

  console.log(chalk.green("Response:"));
  console.log(response);
  console.log();
};

main();
