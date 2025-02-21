import mongoose from "mongoose";

const conn = mongoose.createConnection(
  process.env.MONGO_CONNECTION_STRING || "",
);

conn.useDb("");

const Schema = mongoose.Schema;

const UserContextSchema = new Schema({
  profileId: {
    type: Schema.Types.ObjectId,
    ref: "profiles",
    required: true,
    unique: true,
  },
  timezone: { type: String, default: "UTC" },
  locale: { type: String, default: "en-US" },
  preferences: {
    tone: String,
  },
  context: String,
});

export const UserContext = conn.model("user_contexts", UserContextSchema);
