const { ApolloServer } = require("apollo-server");
const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require("@apollo/gateway");

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: "core", url: "http://localhost:8080/graphql" },
      { name: "timetrackings", url: "http://localhost:8081/graphql" },
      { name: "analytics", url: "http://localhost:8082/graphql" },
    ],
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
