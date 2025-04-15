import * as fs from "fs";
import { createServer as createHttpServer, IncomingMessage, Server } from "http";
import { createServer as createHttpsServer } from "https";
import { URL } from "url";
import { WebSocket, WebSocketServer } from "ws";

import { handleAudioChat } from "./handlers/audioChatHandler";
import { setupAuthentication } from "./handlers/authHandler";
import { handleTextChat } from "./handlers/textChatHandler";
import { WebSocketContext, WebSocketHandler } from "./types";

export const initializeWebSocketServer = (): void => {
  const PORT = 8098;
  const endpointHandlers: Record<string, WebSocketHandler> = {
    "/chat/text": handleTextChat,
    "/chat/audio": handleAudioChat,
  };

  // Try to load SSL certificates
  let server: Server;
  let protocol = "ws";
  try {
    const certPath = process.env.SSL_CERT_PATH;
    const keyPath = process.env.SSL_KEY_PATH;
    if (!certPath || !keyPath) {
      throw new Error("SSL certificate paths not provided in environment variables");
    }

    // Check if certificate files exist
    if (fs.existsSync(certPath) && fs.existsSync(keyPath)) {
      const sslOptions = {
        cert: fs.readFileSync(certPath),
        key: fs.readFileSync(keyPath),
      };
      server = createHttpsServer(sslOptions);
      protocol = "wss";
      console.log("🔒 SSL certificates loaded successfully");
    } else {
      throw new Error("SSL certificate files not found");
    }
  } catch (error) {
    console.log(`⚠️  SSL setup failed: ${error.message}. Falling back to non-secure WebSocket.`);
    server = createHttpServer();
  }

  // Create WebSocket server attached to HTTP/HTTPS server
  const wss = new WebSocketServer({ server });

  // Log the appropriate protocol
  console.log(`🚀 WebSocket Server ready at ${protocol}://localhost:${PORT}`);

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

  wss.on("close", (code, reason) => {
    console.log(`🔌 WebSocket closed | code: ${code} | reason: ${reason.toString()}`);
  });

  wss.on("error", (err) => {
    console.error("❌ WebSocket error:", err);
  });

  // Start the server
  server.listen(PORT);
};
