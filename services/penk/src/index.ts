import "./bootstrap";

import { ApolloServer } from "apollo-server";

import { ResolverContext, schema } from "./graphql";
import { getProfileByEmail } from "./utils/db/utils";
import { decodeFirebaseJwt } from "./utils/firebase";
import { Profile } from "./utils/types";

const server = new ApolloServer({
  cors: {
    allowedHeaders: "Authorization",
  },
  context: async ({ req }) => {
    const token = req.headers.authorization ? `${req.headers.authorization}`.split(" ")[1] : "";

    let profile: Profile | undefined = undefined;

    if (token) {
      const decodedToken = await decodeFirebaseJwt(token);
      if (!decodedToken?.email) throw new Error("invalid jwt");
      const mongoProfile = await getProfileByEmail(decodedToken.email);
      if (!mongoProfile) throw new Error("profile not found");
      profile = {
        id: mongoProfile.id,
        name: mongoProfile.name,
        email: decodedToken.email,
        currentCharacterId: mongoProfile.current_character_id.toString(),
      };
    }

    return {
      token,
      profile,
    } satisfies ResolverContext;
  },
  schema,
});

server.listen(8099).then(({ url }) => {
  console.log(`🚀 Server ready at ${url}`);
});
