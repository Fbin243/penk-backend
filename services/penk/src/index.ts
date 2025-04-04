import "./bootstrap";

import { startGraphQLServer } from "./services/graphql";
import { startOAuthServer } from "./services/restful";
import { initializeWebSocketServer } from "./services/websocket";

initializeWebSocketServer();
startOAuthServer();
startGraphQLServer();
