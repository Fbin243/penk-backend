import { IncomingMessage } from "http";
import { URL } from "url";
import { WebSocket, WebSocketServer } from "ws";

import { PenKMessageModel } from "../../utils/database/mongo";
import { getPenKData, getPenKMessages, getProfileByEmail } from "../../utils/database/utils";
import { decodeFirebaseJwt } from "../../utils/firebase";
import {
  Message,
  MessageType,
  Ws_InfoType,
  Ws_Message,
  Ws_MessageType,
  Ws_TextChatInput,
} from "../../utils/types";
import { textChatStream } from "../ai/utils";

interface WebSocketContext {
  isAuthenticated: boolean;
  userId: string;
  profileId: string;
  email: string;
}

export const initializeWebSocketServer = () => {
  const wss = new WebSocketServer({ port: 8098 });
  console.log("🚀 WebSocket Server ready at ws://localhost:8098");

  wss.on("connection", async (ws: WebSocket, req: IncomingMessage) => {
    const context: WebSocketContext = {
      isAuthenticated: false,
      userId: "",
      profileId: "",
      email: "",
    };

    // Parse URL parameters for the endpoint
    const url = new URL(req.url || "", `http://${req.headers.host}`);
    const endpoint = url.pathname;

    // Set a timeout for authentication
    const authTimeout = setTimeout(() => {
      if (!context.isAuthenticated) {
        const response: Ws_Message = {
          type: Ws_MessageType.Info,
          content: Ws_InfoType.AuthenticationTimeout,
          timestamp: new Date().toISOString(),
        };
        ws.send(JSON.stringify(response));
        ws.close();
      }
    }, 10000);

    // Handle authentication through messages
    const handleAuthentication = async (message: string) => {
      try {
        const data = JSON.parse(message.toString());

        if (data.type === Ws_MessageType.Auth) {
          const token = data.content;

          try {
            const decodedToken = await decodeFirebaseJwt(token);
            if (decodedToken?.email) {
              const mongoProfile = await getProfileByEmail(decodedToken.email);
              if (mongoProfile) {
                context.isAuthenticated = true;
                context.email = decodedToken.email;
                context.userId = mongoProfile._id.toString();
                context.profileId = mongoProfile.current_character_id.toString();

                // Clear the auth timeout
                clearTimeout(authTimeout);

                // Send successful authentication message
                const response: Ws_Message = {
                  type: Ws_MessageType.Info,
                  content: Ws_InfoType.AuthenticationSuccess,
                  timestamp: new Date().toISOString(),
                };
                ws.send(JSON.stringify(response));

                // Set up the appropriate handlers based on endpoint
                switch (endpoint) {
                  case "/chat/text":
                    handleTextChat(ws, context);
                    break;
                  case "/chat/audio":
                    handleAudioChat(ws, context);
                    break;
                  default: {
                    const response: Ws_Message = {
                      type: Ws_MessageType.Info,
                      content: "Invalid endpoint",
                      timestamp: new Date().toISOString(),
                    };
                    ws.send(JSON.stringify(response));
                    ws.close();
                    break;
                  }
                }

                // Remove the authentication listener after successful auth
                ws.removeListener("message", handleAuthentication);
              } else {
                sendAuthFailure(ws, "Profile not found");
              }
            } else {
              sendAuthFailure(ws, "Invalid token");
            }
          } catch {
            sendAuthFailure(ws, "Token verification failed");
          }
        } else {
          // If the first message is not an auth message, request authentication
          const response: Ws_Message = {
            type: Ws_MessageType.Info,
            content: Ws_InfoType.AuthenticationRequired,
            timestamp: new Date().toISOString(),
          };
          ws.send(JSON.stringify(response));
        }
      } catch (err) {
        console.error("Error processing authentication message:", err);
        sendAuthFailure(ws, "Invalid message format");
      }
    };

    // Add the authentication message handler
    ws.on("message", handleAuthentication);

    // Helper function to send auth failure responses
    function sendAuthFailure(ws: WebSocket, reason: string) {
      const response: Ws_Message = {
        type: Ws_MessageType.Info,
        content: Ws_InfoType.AuthenticationFailed,
        timestamp: new Date().toISOString(),
      };
      ws.send(JSON.stringify(response));
      console.error(`Authentication failed: ${reason}`);
    }
  });
};

const handleTextChat = (ws: WebSocket, context: WebSocketContext) => {
  ws.on("message", async (message: string) => {
    try {
      const data = JSON.parse(message.toString()) as Ws_TextChatInput;

      const [penkData, penkMessages] = await Promise.all([
        getPenKData(context.userId),
        getPenKMessages(context.profileId),
      ]);

      const userMessage: Message = {
        type: MessageType.UserMessage,
        content: data.content,
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
            newMessage: data.content,
          },
          // This callback is called for each chunk of the response
          (chunkContent, timestamp) => {
            const chunkResponse: Ws_Message = {
              type: Ws_MessageType.TextChat,
              content: chunkContent,
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
                content: Ws_InfoType.MessageStreamCompleted,
                timestamp: new Date().toISOString(),
              };
              ws.send(JSON.stringify(completionResponse));
            }
          },
        );
      } catch (error) {
        console.error("Error in streaming chat:", error);
      }
    } catch (error) {
      console.error("Error processing message:", error);
      const response: Ws_Message = {
        type: Ws_MessageType.Info,
        content: "Error processing message",
        timestamp: new Date().toISOString(),
      };
      ws.send(JSON.stringify(response));
    }
  });

  ws.on("close", () => {
    console.log("Text chat connection closed");
  });
};

// This is a placeholder implementation - will be implemented later
const handleAudioChat = (ws: WebSocket, context: WebSocketContext) => {
  // Using underscore prefix for unused parameters to avoid linter errors
  console.log(`Audio chat connection established for user: ${context.email}`);

  ws.on("message", () => {
    const response: Ws_Message = {
      type: Ws_MessageType.Info,
      content: "Audio chat not implemented yet",
      timestamp: new Date().toISOString(),
    };
    ws.send(JSON.stringify(response));
  });

  ws.on("close", () => {
    console.log("Audio chat connection closed");
  });
};
