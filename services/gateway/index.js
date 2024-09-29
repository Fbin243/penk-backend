const { ApolloServer } = require("apollo-server");
const { ApolloGateway, IntrospectAndCompose } = require("@apollo/gateway");

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: "core", url: "http://localhost:8080/graphql" },
      { name: "timetrackings", url: "http://localhost:8081/graphql" },
      { name: "analytics", url: "http://localhost:8082/graphql" },
    ],
  }),
});

const server = new ApolloServer({
  gateway,

  subscriptions: false,
});

server.listen({ port: 8070 }).then(({ url }) => {
  console.log(`🚀 Server ready at ${url}`);
});
