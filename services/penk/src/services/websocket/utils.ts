import { WebSocket } from "ws";

import { Ws_InfoType, Ws_Message, Ws_MessageType } from "../../utils/types";

export const sendInfoResponse = (ws: WebSocket, infoType: Ws_InfoType): void => {
  const response: Ws_Message = {
    type: Ws_MessageType.Info,
    data: infoType,
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(response));
};

export const sendErrorResponse = (ws: WebSocket, message: string): void => {
  const response: Ws_Message = {
    type: Ws_MessageType.Info,
    data: message,
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(response));
};
