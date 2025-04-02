// import chalk from "chalk";
import mongoose from "mongoose";

import { PenKContextModel, PenKMessageModel, ProfileModel } from "./mongo";

export const getProfileByEmail = async (email: string) => {
  const profile = await ProfileModel.findOne({ email });
  return profile;
};

export const getPenKMessages = async (profileId: string) => {
  const messages = await PenKMessageModel.find({ profile_id: profileId });
  return messages;
};

const convertObjectsToStrings = (result: object): object => {
  const convertedResult = {};

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
        from: "metrics",
        localField: "profile.current_character_id",
        foreignField: "character_id",
        as: "metrics",
      },
    },
    {
      $lookup: {
        from: "habits",
        localField: "profile.current_character_id",
        foreignField: "character_id",
        as: "habits",
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
        metrics: {
          $map: {
            input: "$metrics",
            as: "metric",
            in: {
              id: "$$metric._id",
              name: "$$metric.name",
              value: "$$metric.value",
              unit: "$$metric.unit",
              category_id: "$$metric.category_id",
            },
          },
        },
        habits: {
          $map: {
            input: "$habits",
            as: "habit",
            in: {
              id: "$$habit._id",
              category_id: "$$habit.category_id",
              name: "$$habit.name",
              value: "$$habit.value",
              unit: "$$habit.unit",
              completion_type: "$$habit.completion_type",
              start_time: "$$habit.start_time",
              end_time: "$$habit.end_time",
              frequency: "$$habit.frequency",
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

  // console.log(chalk.green("[PenK Context]"));
  // console.dir(userData, { depth: null, colors: true });
  // console.log();

  return userData;
};
