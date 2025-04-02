import { WebSocket } from "ws";

import { PenKMessageModel } from "../../../utils/database/mongo";
import { getPenKData, getPenKMessages } from "../../../utils/database/utils";
import {
  Message,
  MessageType,
  Ws_InfoType,
  Ws_Message,
  Ws_MessageType,
} from "../../../utils/types";
import { textChatStream } from "../../ai/utils";
import { WebSocketContext } from "../types";

export const handleTextChat = (ws: WebSocket, context: WebSocketContext) => {
  console.log(`Text chat connection established for user: ${context.email}`);

  ws.on("message", async (message: string) => {
    try {
      const { data } = JSON.parse(message.toString()) as Ws_Message;

      const [penkData, penkMessages] = await Promise.all([
        getPenKData(context.userId),
        getPenKMessages(context.profileId),
      ]);

      const userMessage: Message = {
        type: MessageType.UserMessage,
        content: data,
        timestamp: new Date().toISOString(),
      };

      try {
        let completeAiMessage: Message | null = null;

        await textChatStream(
          {
            userData: JSON.stringify(penkData),
            history: penkMessages.map((message) => ({
              type: message.type,
              content: message.content,
              timestamp: message.timestamp.toISOString(),
            })),
            newMessage: data,
          },
          // This callback is called for each chunk of the response
          (chunkContent, timestamp) => {
            const chunkResponse: Ws_Message = {
              type: Ws_MessageType.TextStream,
              data: chunkContent,
              timestamp,
            };
            ws.send(JSON.stringify(chunkResponse));
          },
          // This callback is called when streaming is complete
          (aiMessage) => {
            completeAiMessage = aiMessage;

            // Save messages to database after streaming completes
            if (completeAiMessage) {
              PenKMessageModel.insertMany([
                { ...userMessage, profile_id: context.profileId },
                { ...completeAiMessage, profile_id: context.profileId },
              ]);

              // Send a completion notification
              const completionResponse: Ws_Message = {
                type: Ws_MessageType.Info,
                data: Ws_InfoType.MessageStreamCompleted,
                timestamp: new Date().toISOString(),
              };
              ws.send(JSON.stringify(completionResponse));
            }
          },
        );
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

function sendErrorResponse(ws: WebSocket, message: string): void {
  const response: Ws_Message = {
    type: Ws_MessageType.Info,
    data: message,
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(response));
}
