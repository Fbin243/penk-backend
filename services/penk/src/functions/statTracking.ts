import { Metadata } from "@grpc/grpc-js";
import { FunctionDefinition } from "openai/resources";

import { timeTrackingClient } from "../utils/grpc";
import { TimeTracking, TimeTrackingWithFish } from "../utils/types";
import { PenKFunction } from "./types";

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
  name: PenKFunction.createTimeTracking,
  description:
    "Create a focus session for a character on a specific category or no categories. A category must belong to the character.",
  strict: true,
  parameters: {
    type: "object",
    properties: {
      characterId: {
        type: ["string", "null"],
        description: "if no characters inferred, use the current character id",
      },
      categoryId: {
        type: ["string", "null"],
        description: "if no categories inferred, use null",
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
  name: PenKFunction.updateTimeTracking,
  description: "End the current focus session. Return session data, normal fish, gold fish.",
  strict: true,
  parameters: {
    type: "object",
    properties: {},
    required: [],
    additionalProperties: false,
  },
};
