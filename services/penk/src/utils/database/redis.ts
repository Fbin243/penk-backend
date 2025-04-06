import { createClient, RedisClientType } from "redis";

let client: RedisClientType | null = null;

export const getRedisClient = async () => {
  if (!client) {
    client = createClient({
      url: `redis://default:${process.env.REDIS_PASSWORD}@${process.env.REDIS_URI}`,
    });
    await client.connect();
    client.on("error", (err) => {
      console.error("[Redis] Error", err);
    });
  }
  return client;
};
