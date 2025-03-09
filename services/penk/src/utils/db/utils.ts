import mongoose from "mongoose";

import { MessageModel, ProfileModel, UserContextModel } from "./mongo";

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
    id: userContext._id.toString(),
    profileID: userContext.profile_id.toString(),
    timezone: userContext.timezone,
    locale: userContext.locale,
    context: userContext.context || "",
    preferences: {
      tone: userContext.preferences?.tone || "",
    },
  };
};

export const getUserData = async (profileId: string) => {
  const userData = await UserContextModel.aggregate([
    { $match: { profile_id: new mongoose.Types.ObjectId(profileId) } },
    { $project: { _id: 0, profile_id: 1 } },
    {
      $lookup: {
        from: "profiles",
        localField: "profile_id",
        foreignField: "_id",
        as: "profile",
        pipeline: [{ $project: { current_character_id: 1 } }],
      },
    },
    { $unwind: "$profile" },
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

  // console.log("--> user data");
  // console.dir(userData, { depth: null, colors: true });

  return userData;
};

export const getMessages = async (profileId: string) => {
  const messages = await MessageModel.find({
    profile_id: profileId,
  });

  return messages;
};
