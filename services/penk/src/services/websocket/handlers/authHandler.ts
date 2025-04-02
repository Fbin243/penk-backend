import { WebSocket } from "ws";

import { getProfileByEmail } from "../../../utils/database/utils";
import { decodeFirebaseJwt } from "../../../utils/firebase";
import { Ws_InfoType, Ws_Message, Ws_MessageType } from "../../../utils/types";
import { WebSocketContext, WebSocketHandler } from "../types";

export const setupAuthentication = (
  ws: WebSocket,
  context: WebSocketContext,
  endpoint: string,
  endpointHandlers: Record<string, WebSocketHandler>,
  authTimeoutMs = 10000,
): void => {
  // Set a timeout for authentication
  const authTimeout = setTimeout(() => {
    if (!context.isAuthenticated) {
      const response: Ws_Message = {
        type: Ws_MessageType.Info,
        data: Ws_InfoType.AuthenticationTimeout,
        timestamp: new Date().toISOString(),
      };
      ws.send(JSON.stringify(response));
      ws.close();
    }
  }, authTimeoutMs);

  // Handle authentication through messages
  const handleAuthentication = async (message: string) => {
    try {
      const { data, type } = JSON.parse(message.toString()) as Ws_Message;

      if (type === Ws_MessageType.Auth) {
        const token = data;

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
                data: Ws_InfoType.AuthenticationSuccess,
                timestamp: new Date().toISOString(),
              };
              ws.send(JSON.stringify(response));

              // Set up the appropriate handlers based on endpoint
              const handler = endpointHandlers[endpoint];
              if (handler) {
                handler(ws, context);
              } else {
                sendInvalidEndpointResponse(ws);
              }

              // Remove the authentication listener after successful auth
              ws.removeListener("message", handleAuthentication);
            } else {
              sendAuthFailure(ws, "Profile not found");
            }
          } else {
            sendAuthFailure(ws, "Invalid token");
          }
        } catch (authError) {
          console.error("Token verification failed:", authError);
          sendAuthFailure(ws, "Token verification failed");
        }
      } else {
        const response: Ws_Message = {
          type: Ws_MessageType.Info,
          data: Ws_InfoType.AuthenticationRequired,
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
};

function sendInvalidEndpointResponse(ws: WebSocket): void {
  const response: Ws_Message = {
    type: Ws_MessageType.Info,
    data: "Invalid endpoint",
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(response));
  ws.close();
}

function sendAuthFailure(ws: WebSocket, reason: string): void {
  const response: Ws_Message = {
    type: Ws_MessageType.Info,
    data: Ws_InfoType.AuthenticationFailed,
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(response));
  console.error(`Authentication failed: ${reason}`);
}
