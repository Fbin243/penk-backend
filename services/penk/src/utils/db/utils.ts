import mongoose from "mongoose";

import { MessageModel, ProfileModel, UserContextModel } from "./mongo";

const convertObjectIdsToStrings = (result: object): object => {
  const convertedResult = {};

  for (const key in result) {
    if (Object.prototype.hasOwnProperty.call(result, key)) {
      const value = result[key];

      if (Array.isArray(value)) {
        // Handle arrays recursively
        convertedResult[key] = value.map((item) => convertObjectIdsToStrings(item));
      } else if (
        typeof value === "object" &&
        value !== null &&
        value.constructor.name === "ObjectId"
      ) {
        // Check for ObjectId and convert to string
        convertedResult[key] = value.toString();
      } else if (typeof value === "object" && value !== null) {
        // handle nested objects recursively
        convertedResult[key] = convertObjectIdsToStrings(value);
      } else {
        // Keep other values as they are
        convertedResult[key] = value;
      }
    }
  }

  return convertedResult;
};

export const getProfileByEmail = async (email: string) => {
  const profile = await ProfileModel.findOne({
    email,
  });
  return profile;
};

export const getUserContext = async (profileId: string) => {
  let userContext = await UserContextModel.findOne({
    profile_id: profileId,
  });

  if (!userContext) {
    userContext = await UserContextModel.create({
      profile_id: profileId,
      locale: "vi",
      timezone: "Asia/Ho_Chi_Minh",
      context: "",
      preferences: {
        tone: "funny",
      },
    });
  }

  return {
    timezone: userContext.timezone,
    locale: userContext.locale,
    context: userContext.context || "",
    preferences: {
      tone: userContext.preferences?.tone || "",
    },
  };
};

export const getUserData = async (profileId: string) => {
  const userContext = await getUserContext(profileId);

  const aggregatedData = await ProfileModel.aggregate([
    { $match: { _id: new mongoose.Types.ObjectId(profileId) } },
    {
      $project: {
        _id: 0,
        profile_id: "$_id",
        current_character_id: 1,
      },
    },
    {
      $lookup: {
        from: "characters",
        localField: "profile_id",
        foreignField: "profile_id",
        as: "characters",
        pipeline: [
          {
            $project: {
              name: 1,
              categories: {
                $map: {
                  input: "$categories",
                  as: "category",
                  in: {
                    _id: "$$category._id",
                    name: "$$category.name",
                    description: "$$category.description",
                    metrics: {
                      $map: {
                        input: "$$category.metrics",
                        as: "metric",
                        in: {
                          _id: "$$metric._id",
                          name: "$$metric.name",
                          value: "$$metric.value",
                          unit: "$$metric.unit",
                        },
                      },
                    },
                  },
                },
              },
              metrics: {
                $map: {
                  input: "$metrics",
                  as: "metric",
                  in: {
                    _id: "$$metric._id",
                    name: "$$metric.name",
                    value: "$$metric.value",
                    unit: "$$metric.unit",
                  },
                },
              },
            },
          },
          {
            $lookup: {
              from: "goals",
              localField: "_id",
              foreignField: "character_id",
              as: "goals",
              pipeline: [
                {
                  $project: {
                    name: 1,
                    description: 1,
                    start_date: 1,
                    end_date: 1,
                    status: 1,
                    target: 1,
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ]);

  const userData = { ...convertObjectIdsToStrings(aggregatedData[0]), context: userContext };

  // console.log("[User Data]");
  // console.dir(userData, { depth: null, colors: true });
  // console.log();

  return userData;
};

export const getMessages = async (profileId: string) => {
  const messages = await MessageModel.find({
    profile_id: profileId,
  });

  return messages;
};
