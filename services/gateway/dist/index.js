"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const gateway_1 = require("@apollo/gateway");
const apollo_server_1 = require("apollo-server");
const node_fetch_1 = __importDefault(require("node-fetch"));
require("dotenv").config();
const subgraphs = [
    { name: "core", url: process.env.CORE_URL || "http://localhost:8080/graphql" },
    { name: "analytic", url: process.env.ANALYTIC_URL || "http://localhost:8082/graphql" },
    { name: "notification", url: process.env.NOTIFICATION_URL || "http://localhost:8084/graphql" },
    { name: "currency", url: process.env.CURRENCY_URL || "http://localhost:8085/graphql" },
];
async function startGateway() {
    // Check if all subgraphs are reachable
    const reachableSubgraphs = [];
    await Promise.all(subgraphs.map(async ({ name, url }) => {
        try {
            const response = await (0, node_fetch_1.default)(url);
            if (response.ok) {
                console.log(`Subgraph at ${url} is reachable`);
            }
            reachableSubgraphs.push({ name, url });
        }
        catch (e) {
            console.error(`Subgraph at ${url} is not reachable`);
        }
    }));
    // Create the apollo gateway
    const gateway = new gateway_1.ApolloGateway({
        supergraphSdl: new gateway_1.IntrospectAndCompose({
            subgraphs: reachableSubgraphs,
        }),
        buildService({ name, url }) {
            return new gateway_1.RemoteGraphQLDataSource({
                url,
                willSendRequest({ request, context }) {
                    request.http?.headers.append("Authorization", context.token);
                    request.http?.headers.append("X-Device-Id", context.deviceId);
                },
            });
        },
    });
    // Start the server with the gateway
    const server = new apollo_server_1.ApolloServer({
        gateway,
        context: ({ req }) => {
            return {
                token: req.headers.authorization,
                deviceId: req.headers["x-device-id"],
            };
        },
    });
    server.listen({ port: 8070 }).then(({ url }) => {
        console.log(`🚀 Server ready at ${url}`);
    });
}
startGateway();
