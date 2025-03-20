import { FunctionCallType } from "../utils/types";
import { createTimeTracking, updateTimeTracking } from "./statTracking";

export const openaiPenKMap = {
  [FunctionCallType.CreateTimeTracking]: createTimeTracking,
  [FunctionCallType.UpdateTimeTracking]: updateTimeTracking,
};
