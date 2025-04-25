/* eslint-disable @typescript-eslint/no-explicit-any */
import mongoose from "mongoose";

import { decrypt } from "../encrypt";
import { refreshToken } from "../googleapis/auth";
import { LinkedAccount } from "../types";
import { OAuthTokenModel, PenKContextModel, PenKMessageModel, ProfileModel } from "./mongo";
import { getRedisClient } from "./redis";

export const getProfileByEmail = async (email: string) => {
  const profile = await ProfileModel.findOne({ email });
  return profile;
};

export const getPenKMessages = async (
  profileId: string,
  offset: number = 0,
  limit: number = 20,
) => {
  const query = PenKMessageModel.find({ profile_id: profileId })
    .sort({ timestamp: -1 })
    .skip(offset)
    .limit(limit);
  const messages = await query;
  return messages;
};

const convertObjectsToStrings = (result: any): any => {
  const convertedResult: any = {};

  for (const key in result) {
    if (Object.prototype.hasOwnProperty.call(result, key)) {
      const value = result[key];

      if (Array.isArray(value)) {
        // Handle arrays recursively
        convertedResult[key] = value.map((item) => convertObjectsToStrings(item));
      } else if (
        typeof value === "object" &&
        value !== null &&
        value.constructor.name === "ObjectId"
      ) {
        // Check for ObjectId and convert to string
        convertedResult[key] = value.toString();
      } else if (typeof value === "object" && value !== null && value.constructor.name === "Date") {
        // Convert Date to ISO string
        convertedResult[key] = value.toISOString();
      } else if (typeof value === "object" && value !== null) {
        // handle nested objects recursively
        convertedResult[key] = convertObjectsToStrings(value);
      } else {
        // Keep other values as they are
        convertedResult[key] = value;
      }
    }
  }

  return convertedResult;
};

export const getPenKContext = async (userId: string) => {
  const penkContext = await PenKContextModel.findOne({ user_id: userId });
  return penkContext;
};

export const getPenKData = async (userId: string) => {
  const aggregatedData = await PenKContextModel.aggregate([
    { $match: { user_id: new mongoose.Types.ObjectId(userId) } },
    {
      $lookup: {
        from: "profiles",
        localField: "user_id",
        foreignField: "_id",
        as: "profile",
      },
    },
    { $unwind: "$profile" },
    {
      $lookup: {
        from: "categories",
        localField: "profile.current_character_id",
        foreignField: "character_id",
        as: "categories",
      },
    },
    {
      $lookup: {
        from: "goals",
        localField: "profile.current_character_id",
        foreignField: "character_id",
        as: "goals",
      },
    },
    {
      $project: {
        _id: 0,
        timezone: 1,
        locale: 1,
        context: 1,
        firebase_uid: "$profile.firebase_uid",
        profile_id: "$profile.current_character_id",
        categories: {
          $map: {
            input: "$categories",
            as: "category",
            in: {
              id: "$$category._id",
              name: "$$category.name",
              description: "$$category.description",
            },
          },
        },
        goals: {
          $map: {
            input: "$goals",
            as: "goal",
            in: {
              id: "$$goal._id",
              name: "$$goal.name",
              description: "$$goal.description",
              start_time: "$$goal.start_time",
              end_time: "$$goal.end_time",
              completion_time: "$$goal.completion_time",
              metrics: {
                $map: {
                  input: "$$goal.metrics",
                  as: "metric",
                  in: {
                    id: "$$metric.id",
                    condition: "$$metric.condition",
                    target_value: "$$metric.target_value",
                    range_value: "$$metric.range_value",
                  },
                },
              },
              checkboxes: {
                $map: {
                  input: "$$goal.checkboxes",
                  as: "checkbox",
                  in: {
                    id: "$$checkbox.id",
                    name: "$$checkbox.name",
                    value: "$$checkbox.value",
                  },
                },
              },
            },
          },
        },
      },
    },
  ]);

  const userData = convertObjectsToStrings(aggregatedData[0]);

  return userData;
};

export const getLinkedAccounts = async (profileId: string) => {
  const redisClient = await getRedisClient();

  const cached = await redisClient.get(`linked_accounts:${profileId}`);
  if (cached) {
    return JSON.parse(cached) as LinkedAccount[];
  }

  const tokens = await OAuthTokenModel.find({ profile_id: profileId });
  const linkedAccounts: LinkedAccount[] = [];
  for (const token of tokens) {
    try {
      const accessToken = await refreshToken(decrypt(token.refresh_token));
      if (accessToken) {
        linkedAccounts.push({
          id: token._id.toString(),
          email: token.email,
          type: token.type,
          accessToken,
        });
      }
    } catch (error) {
      if ((error as any)?.response?.data?.error === "invalid_grant") {
        await OAuthTokenModel.deleteOne({ _id: token._id });
      } else {
        throw error;
      }
    }
  }

  redisClient.set(`linked_accounts:${profileId}`, JSON.stringify(linkedAccounts), {
    EX: 60 * 50, // 50 minutes
  });

  return linkedAccounts;
};
