import { WebSocket } from "ws";

import { textStream } from "../../../utils/ai";
import { MembershipModel, PenKMessageModel } from "../../../utils/database/mongo";
import { getPenKData, getPenKMessages } from "../../../utils/database/utils";
import { Ws_Message, Ws_MessageType } from "../../../utils/types";
import { WebSocketContext } from "../types";
import { sendErrorResponse } from "../utils";

const costMultiplier = 5;

export const handleTextChat = (ws: WebSocket, context: WebSocketContext) => {
  console.log(`Text chat connection established for user: ${context.email}`);

  ws.on("message", async (message: string) => {
    try {
      const { data } = JSON.parse(message.toString()) as Ws_Message;

      const membership = await MembershipModel.findOne({ email: context.email }).catch((error) => {
        console.error("Error fetching membership data:", error);
        return null;
      });
      const monthlyCredit = membership?.monthly_credit || 0;
      const persistentCredit = membership?.persistent_credit || 0;
      if (monthlyCredit <= 0 && persistentCredit <= 0) {
        sendErrorResponse(ws, "Not enough credits");
        return;
      }

      const [penkData, history] = await Promise.all([
        getPenKData(context.userId),
        getPenKMessages(context.profileId),
      ]);

      try {
        const { newPenKMessages, cost } = await textStream(
          {
            userData: JSON.stringify(penkData),
            history: history.reverse().map((message) => ({
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

        const creditCost = cost * costMultiplier;
        const monthlyCreditUsed = monthlyCredit > creditCost ? creditCost : monthlyCredit;
        const persistentCreditUsed = monthlyCredit > creditCost ? 0 : creditCost - monthlyCredit;
        await MembershipModel.findOneAndUpdate(
          { email: context.email },
          {
            $inc: {
              monthly_credit: -monthlyCreditUsed,
              persistent_credit: -persistentCreditUsed,
            },
          },
          { new: true },
        ).catch((error) => {
          console.error("Error updating payment data:", error);
          return null;
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
