import mongoose from "mongoose";

import { MessageModel, UserContextModel } from "./mongo";

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

  const userData = await UserContextModel.aggregate([
    {
      $match: { profile_id: new mongoose.Types.ObjectId(profileId) },
    },
    {
      $lookup: {
        from: "profiles",
        localField: "profile_id",
        foreignField: "_id",
        as: "profile",
        pipeline: [
          {
            $project: {
              _id: 0,
              name: 1,
              current_character_id: 1,
            },
          },
        ],
      },
    },
    {
      $unwind: "$profile",
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
              // custom_metrics: 1, // TODO: Replace with "categories"
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

  console.dir(userData, { depth: null, colors: true });
  console.log(JSON.stringify(userData));

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

export const getMessages = async (profileId: string) => {
  const messages = await MessageModel.find({
    profile_id: profileId,
  });

  return messages;
};
