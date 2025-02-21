import { createClient, RedisClientType } from "redis";

let client: RedisClientType | null = null;

export const getRedisClient = async () => {
  if (!client) {
    client = createClient({
      url: process.env.REDIS_CONNECTION_STRING,
    });

    client.on("error", (err) => console.log("Redis Client Error", err));

    await client.connect();
  }

  return client;
};
