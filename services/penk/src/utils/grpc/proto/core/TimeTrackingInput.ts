// Original file: ../../proto/core/timetracking_message.proto

import type {
  EntityType as _core_EntityType,
  EntityType__Output as _core_EntityType__Output,
} from "./EntityType";
import type { Long } from "@grpc/proto-loader";

export interface TimeTrackingInput {
  referenceId?: string;
  referenceType?: _core_EntityType;
  timestamp?: number | string | Long;
  duration?: number | string | Long;
  _referenceId?: "referenceId";
  _referenceType?: "referenceType";
}

export interface TimeTrackingInput__Output {
  referenceId?: string;
  referenceType?: _core_EntityType__Output;
  timestamp: string;
  duration: string;
  _referenceId: "referenceId";
  _referenceType: "referenceType";
}
