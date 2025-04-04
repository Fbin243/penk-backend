import mongoose from "mongoose";

import { LinkedAccountType, MessageType } from "../types";

const conn = mongoose.createConnection(
  `mongodb+srv://${process.env.MONGO_USER}:${process.env.MONGO_PASSWORD}@${process.env.MONGO_ADDRESS}/${process.env.MONGO_DATABASE_NAME}?retryWrites=true&w=majority`,
  {
    autoIndex: true,
  },
);

const Schema = mongoose.Schema;

const ProfileSchema = new Schema({
  name: { type: String, required: true },
  email: { type: String, required: true },
  current_character_id: { type: Schema.Types.ObjectId, required: true },
});

export const ProfileModel = conn.model("profiles", ProfileSchema);

const PenKContextSchema = new Schema({
  user_id: {
    type: Schema.Types.ObjectId,
    ref: "profiles",
    required: true,
    unique: true,
  },
  timezone: { type: String, default: "Asia/Saigon" },
  locale: { type: String, default: "vi-VN" },
  context: { type: String, default: "" },
});

export const PenKContextModel = conn.model("penk_contexts", PenKContextSchema);

const PenKMessageSchema = new Schema(
  {
    profile_id: {
      type: Schema.Types.ObjectId,
      ref: "characters",
      required: true,
    },
    timestamp: { type: Date, required: true },
    content: { type: String, required: true },
    type: {
      type: String,
      required: true,
      enum: [MessageType.UserMessage, MessageType.AiMessage],
    },
  },
  {
    timeseries: {
      metaField: "profile_id",
      timeField: "timestamp",
      granularity: "seconds",
    },
    versionKey: false,
  },
);

PenKMessageSchema.index(
  { timestamp: 1 },
  { expireAfterSeconds: 60 * 60 * 24, partialFilterExpression: { profile_id: { $exists: true } } },
);

export const PenKMessageModel = conn.model("penk_messages", PenKMessageSchema);

const PenKUsageSchema = new Schema({
  profile_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  total_cost: { type: Number, required: true },
  text_chat_count: { type: Number, required: true },
  voice_chat_count: { type: Number, required: true },
});

export const PenKUsageModel = conn.model("penk_usages", PenKUsageSchema);

// 1 profile can link multiple oauth tokens
export const OAuthTokenSchema = new Schema({
  profile_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  email: { type: String, required: true },
  type: {
    type: String,
    required: true,
    enum: [LinkedAccountType.Gmail, LinkedAccountType.GoogleCalendar],
  },
  refresh_token: { type: String, required: true },
});

export const OAuthTokenModel = conn.model("oauth_tokens", OAuthTokenSchema);
