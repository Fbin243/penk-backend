import { WebSocket } from "ws";

export interface WebSocketContext {
  isAuthenticated: boolean;
  userId: string;
  profileId: string;
  email: string;
}

export type WebSocketHandler = (ws: WebSocket, context: WebSocketContext) => void;
