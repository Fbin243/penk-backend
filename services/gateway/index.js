const { ApolloServer } = require("apollo-server");
const {
  ApolloGateway,
  IntrospectAndCompose,
  RemoteGraphQLDataSource,
} = require("@apollo/gateway");
const fetch = require("node-fetch");

const subgraphs = [
  { name: "core", url: "http://localhost:8080/graphql" },
  { name: "timetracking", url: "http://localhost:8082/graphql" },
  { name: "analytic", url: "http://localhost:8083/graphql" },
  { name: "notification", url: "http://localhost:8084/graphql" },
  { name: "currency", url: "http://localhost:8085/graphql" },
  { name: "penk", url: "http://localhost:8099/graphql" },
];

async function startGateway() {
  // Check if all subgraphs are reachable
  const reachableSubgraphs = [];
  await Promise.all(
    subgraphs.map(async ({ name, url }) => {
      try {
        const response = await fetch(url);
        if (response.ok) {
          console.log(`Subgraph at ${url} is reachable`);
        }

        reachableSubgraphs.push({ name, url });
      } catch (e) {
        console.error(`Subgraph at ${url} is not reachable`);
      }
    }),
  );

  // Create the apollo gateway
  const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
      subgraphs: reachableSubgraphs,
    }),
    buildService({ name, url }) {
      return new RemoteGraphQLDataSource({
        url,
        willSendRequest({ request, context }) {
          request.http.headers.set("Authorization", context.token);
        },
      });
    },
  });

  // Start the server with the gateway
  const server = new ApolloServer({
    gateway,
    context: ({ req }) => {
      return {
        token: req.headers.authorization,
      };
    },
    subscriptions: false,
  });

  server.listen({ port: 8070 }).then(({ url }) => {
    console.log(`🚀 Server ready at ${url}`);
  });
}

startGateway();
