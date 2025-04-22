// Original file: ../../proto/core/timetracking_message.proto

import type {
  EntityType as _core_EntityType,
  EntityType__Output as _core_EntityType__Output,
} from "./EntityType";
import type { Long } from "@grpc/proto-loader";

export interface TimeTracking {
  id?: string;
  characterId?: string;
  categoryId?: string;
  referenceId?: string;
  referenceType?: _core_EntityType;
  timestamp?: number | string | Long;
  duration?: number | string | Long;
  _categoryId?: "categoryId";
  _referenceId?: "referenceId";
  _referenceType?: "referenceType";
}

export interface TimeTracking__Output {
  id: string;
  characterId: string;
  categoryId?: string;
  referenceId?: string;
  referenceType?: _core_EntityType__Output;
  timestamp: string;
  duration: string;
  _categoryId: "categoryId";
  _referenceId: "referenceId";
  _referenceType: "referenceType";
}
