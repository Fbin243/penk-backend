// Original file: ../../proto/core/goal_message.proto

import type {
  MetricCondition as _core_MetricCondition,
  MetricCondition__Output as _core_MetricCondition__Output,
} from "./MetricCondition";
import type {
  RangeInput as _core_RangeInput,
  RangeInput__Output as _core_RangeInput__Output,
} from "./RangeInput";

export interface GoalMetricInput {
  id?: string;
  condition?: _core_MetricCondition;
  targetValue?: number | string;
  rangeValue?: _core_RangeInput | null;
  _targetValue?: "targetValue";
  _rangeValue?: "rangeValue";
}

export interface GoalMetricInput__Output {
  id: string;
  condition: _core_MetricCondition__Output;
  targetValue?: number;
  rangeValue?: _core_RangeInput__Output | null;
  _targetValue: "targetValue";
  _rangeValue: "rangeValue";
}
