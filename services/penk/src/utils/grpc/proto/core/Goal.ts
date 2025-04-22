// Original file: ../../proto/core/goal_message.proto

import type {
  GoalStatus as _core_GoalStatus,
  GoalStatus__Output as _core_GoalStatus__Output,
} from "./GoalStatus";
import type {
  GoalMetric as _core_GoalMetric,
  GoalMetric__Output as _core_GoalMetric__Output,
} from "./GoalMetric";
import type {
  Checkbox as _core_Checkbox,
  Checkbox__Output as _core_Checkbox__Output,
} from "./Checkbox";
import type { Long } from "@grpc/proto-loader";

export interface Goal {
  id?: string;
  createdAt?: number | string | Long;
  updatedAt?: number | string | Long;
  characterId?: string;
  name?: string;
  description?: string;
  startTime?: number | string | Long;
  endTime?: number | string | Long;
  status?: _core_GoalStatus;
  metrics?: _core_GoalMetric[];
  checkboxes?: _core_Checkbox[];
}

export interface Goal__Output {
  id: string;
  createdAt: string;
  updatedAt: string;
  characterId: string;
  name: string;
  description: string;
  startTime: string;
  endTime: string;
  status: _core_GoalStatus__Output;
  metrics: _core_GoalMetric__Output[];
  checkboxes: _core_Checkbox__Output[];
}
