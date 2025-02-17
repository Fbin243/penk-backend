import { createTimeTracking, updateTimeTracking } from "./statTracking";
import { PenKFunction } from "./types";

export const openaiPenKMap = {
  [PenKFunction.createTimeTracking]: createTimeTracking,
  [PenKFunction.updateTimeTracking]: updateTimeTracking,
};
