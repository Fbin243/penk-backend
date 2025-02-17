import { FunctionDefinition } from "openai/resources";
import { PenKFunction, TimeTracking, TimeTrackingWithFish } from "./types";

const localTimeTrackings: TimeTracking[] = [];

export const createTimeTracking = async (props: {
  characterId: string;
  categoryId?: string;
}): Promise<TimeTracking> => {
  const testTimeTracking = {
    id: "tt123",
    startTime: new Date().toISOString(),
    characterId: props.characterId,
    categoryId: props.categoryId,
  };

  localTimeTrackings.push(testTimeTracking);

  return testTimeTracking;
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

export const updateTimeTracking =
  async (props: {}): Promise<TimeTrackingWithFish> => {
    localTimeTrackings[0].endTime = new Date().toISOString();

    return {
      timeTracking: localTimeTrackings[0],
      normal: 10,
      gold: 0,
    };
  };

export const openaiUpdateTimeTracking: FunctionDefinition = {
  name: PenKFunction.updateTimeTracking,
  description:
    "End the current focus session. Return session data, normal fish, gold fish.",
  strict: true,
  parameters: {
    type: "object",
    properties: {},
    required: [],
    additionalProperties: false,
  },
};
