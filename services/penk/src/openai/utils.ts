import { Metadata } from "@grpc/grpc-js";
import { ChatCompletionMessageToolCall, ChatCompletionToolMessageParam } from "openai/resources";

import { openaiPenKMap } from "../functions";

export const handleToolCalls = async (
  toolCalls: ChatCompletionMessageToolCall[],
  metadata: Metadata,
): Promise<ChatCompletionToolMessageParam[]> => {
  const messages: ChatCompletionToolMessageParam[] = [];

  try {
    const toolCallPromises = toolCalls.map(async (toolCall) => {
      try {
        const args = JSON.parse(toolCall.function.arguments);
        console.log(`[Function calling] ${toolCall.function.name}, args:`);
        console.dir(args, { depth: null, colors: true });
        console.log();

        const result = await openaiPenKMap[toolCall.function.name](args, metadata);
        console.log(`[Result injecting] ${toolCall.function.name}`);
        console.dir(result, { depth: null, colors: true });
        console.log();

        return { toolCallId: toolCall.id, result };
      } catch (error) {
        console.error(
          `Error processing tool call ${toolCall.function.name} (${toolCall.id}):`,
          error,
        );
        return {
          toolCallId: toolCall.id,
          result: `Error processing tool call: ${error}`,
        };
      }
    });

    const results = await Promise.all(toolCallPromises);

    results.forEach(({ toolCallId, result }) => {
      messages.push({
        role: "tool",
        tool_call_id: toolCallId,
        content: JSON.stringify(result),
      });
    });
  } catch (error) {
    console.error("Error in handleToolCalls:", error);
  }

  return messages;
};
