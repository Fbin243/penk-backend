import { IncomingMessage } from "http";
import { URL } from "url";
import { WebSocket, WebSocketServer } from "ws";

import { handleAudioChat } from "./handlers/audioChatHandler";
import { setupAuthentication } from "./handlers/authHandler";
import { handleTextChat } from "./handlers/textChatHandler";
import { WebSocketContext, WebSocketHandler } from "./types";

export const initializeWebSocketServer = (): void => {
  const wss = new WebSocketServer({ port: 8098 });
  console.log("🚀 WebSocket Server ready at ws://localhost:8098");

  // Define endpoint handlers
  const endpointHandlers: Record<string, WebSocketHandler> = {
    "/chat/text": handleTextChat,
    "/chat/audio": handleAudioChat,
  };

  wss.on("connection", (ws: WebSocket, req: IncomingMessage) => {
    const context: WebSocketContext = {
      isAuthenticated: false,
      userId: "",
      profileId: "",
      email: "",
    };

    const url = new URL(req.url || "", `http://${req.headers.host}`);
    const endpoint = url.pathname;

    setupAuthentication(ws, context, endpoint, endpointHandlers);
  });
};
