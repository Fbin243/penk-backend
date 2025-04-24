import { WebSocket } from "ws";

import { textStream } from "../../../utils/ai";
import { PenKMessageModel, PenKUsageModel } from "../../../utils/database/mongo";
import { getPenKData, getPenKMessages } from "../../../utils/database/utils";
import { Ws_Message, Ws_MessageType } from "../../../utils/types";
import { WebSocketContext } from "../types";
import { sendErrorResponse } from "../utils";

export const handleTextChat = (ws: WebSocket, context: WebSocketContext) => {
  console.log(`Text chat connection established for user: ${context.email}`);

  ws.on("message", async (message: string) => {
    try {
      const { data } = JSON.parse(message.toString()) as Ws_Message;

      const [penkData, history] = await Promise.all([
        getPenKData(context.userId),
        getPenKMessages(context.profileId),
      ]);

      try {
        const { newPenKMessages, cost } = await textStream(
          {
            userData: JSON.stringify(penkData),
            history: history.map((message) => ({
              type: message.type,
              content: message.content,
              timestamp: message.timestamp.toISOString(),
            })),
            newMessage: data,
          },
          // This callback is called for each chunk of the response
          (chunkContent) => {
            const chunkResponse: Ws_Message = {
              type: Ws_MessageType.TextStream,
              data: chunkContent,
              timestamp: new Date().toISOString(),
            };
            ws.send(JSON.stringify(chunkResponse));
          },
          // This callback is called when a tool is called
          (toolName) => {
            const toolCallResponse: Ws_Message = {
              type: Ws_MessageType.ToolCall,
              data: toolName,
              timestamp: new Date().toISOString(),
            };
            ws.send(JSON.stringify(toolCallResponse));
          },
          // This callback is called when streaming is complete
          (response) => {
            const completionResponse: Ws_Message = {
              type: Ws_MessageType.TextStreamEnded,
              data: response,
              timestamp: new Date().toISOString(),
            };
            ws.send(JSON.stringify(completionResponse));
          },
        );

        if (newPenKMessages.length > 0) {
          PenKMessageModel.insertMany(
            newPenKMessages.map((m) => ({
              ...m,
              profile_id: context.profileId,
            })),
          ).catch((error) => {
            console.error("Error saving messages:", error);
          });
        }

        console.log("Total cost:", cost);
        console.log();
        PenKUsageModel.updateOne(
          { profile_id: context.profileId },
          { $inc: { total_cost: cost, text_chat_count: 1 } },
          { upsert: true },
        ).catch((error) => {
          console.error("Error updating usage:", error);
        });
      } catch (error) {
        console.error("Error in streaming chat:", error);
        sendErrorResponse(ws, "Error in streaming chat");
      }
    } catch (error) {
      console.error("Error processing message:", error);
      sendErrorResponse(ws, "Error processing message");
    }
  });

  ws.on("close", () => {
    console.log("Text chat connection closed");
  });
};
