import { Metadata } from "@grpc/grpc-js";
import { FunctionDefinition } from "openai/resources";

import { timeTrackingClient } from "../utils/grpc";
import { FunctionCallType, TimeTracking, TimeTrackingWithFish } from "../utils/types";

export const createTimeTracking = async (
  props: {
    characterId: string;
    categoryId?: string;
  },
  metadata: Metadata,
): Promise<TimeTracking> => {
  return new Promise<TimeTracking>((resolve, reject) => {
    timeTrackingClient.CreateTimeTracking(
      {
        characterId: props.characterId,
        categoryId: props.categoryId,
        startTime: Date.now() / 1000,
      },
      metadata,
      (err, data) => {
        if (err) {
          reject(err);
          return;
        }

        if (data) {
          const timeTracking: TimeTracking = {
            id: data.id,
            characterId: data.characterId,
            categoryId: data.categoryId,
            startTime: data.startTime,
            endTime: data.endTime,
          };
          resolve(timeTracking);
        } else {
          reject(new Error("Time tracking data not found"));
        }
      },
    );
  });
};

export const openaiCreateTimeTracking: FunctionDefinition = {
  name: FunctionCallType.CreateTimeTracking,
  description:
    "Create a focus session for a character on a specific category id or no category ids. A category id must belong to the character.",
  strict: true,
  parameters: {
    type: "object",
    properties: {
      characterId: {
        type: ["string", "null"],
        description:
          "If no characters inferred, use the current character id. Access this prop via `character._id`.",
      },
      categoryId: {
        type: ["string", "null"],
        description:
          "If no category id is inferred, use null. Access this prop via `character.categories[index]._id`.",
      },
    },
    required: ["characterId", "categoryId"],
    additionalProperties: false,
  },
};

export const updateTimeTracking = async (
  props: object,
  metadata: Metadata,
): Promise<TimeTrackingWithFish> => {
  return new Promise<TimeTrackingWithFish>((resolve, reject) => {
    timeTrackingClient.UpdateTimeTracking(props, metadata, (err, data) => {
      if (err) {
        reject(err);
        return;
      }

      if (data && data.timeTracking) {
        resolve({
          timeTracking: data.timeTracking,
          normal: data.normal,
          gold: data.gold,
        });
      } else {
        reject(new Error("Time tracking data not found"));
      }
    });
  });
};

export const openaiUpdateTimeTracking: FunctionDefinition = {
  name: FunctionCallType.UpdateTimeTracking,
  description: "End the current focus session. Return session data, normal fish, gold fish.",
  strict: true,
  parameters: {
    type: "object",
    properties: {},
    required: [],
    additionalProperties: false,
  },
};
