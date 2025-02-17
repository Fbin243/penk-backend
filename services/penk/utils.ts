import {
  ChatCompletionMessageParam,
  ChatCompletionMessageToolCall,
} from "openai/resources";
import { openaiPenKMap } from "./functions";
import chalk from "chalk";

export const handleToolCalls = async (
  toolCalls: ChatCompletionMessageToolCall[],
): Promise<ChatCompletionMessageParam[]> => {
  const messages: ChatCompletionMessageParam[] = [];

  try {
    const toolCallPromises = toolCalls.map(async (toolCall) => {
      try {
        const args = JSON.parse(toolCall.function.arguments);
        console.log(
          chalk.green(`[Function Detected] ${toolCall.function.name}, args:`),
        );
        console.dir(args, { depth: null, colors: true });
        console.log();

        const result = await openaiPenKMap[toolCall.function.name](args);
        return { toolCallId: toolCall.id, result };
      } catch (error) {
        console.error(`Error processing tool call ${toolCall.id}:`, error);
        return {
          toolCallId: toolCall.id,
          result: "Error processing tool call",
        };
      }
    });

    const results = await Promise.all(toolCallPromises);

    results.forEach(({ toolCallId, result }) => {
      messages.push({
        role: "tool",
        tool_call_id: toolCallId,
        content: `${result}`,
      });
    });
  } catch (error) {
    console.error("Error in handleToolCalls:", error);
  }

  return messages;
};
