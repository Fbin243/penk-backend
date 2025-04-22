// Original file: ../../proto/core/goal_message.proto

import type {
  MetricCondition as _core_MetricCondition,
  MetricCondition__Output as _core_MetricCondition__Output,
} from "./MetricCondition";
import type { Range as _core_Range, Range__Output as _core_Range__Output } from "./Range";

export interface GoalMetric {
  id?: string;
  condition?: _core_MetricCondition;
  targetValue?: number | string;
  rangeValue?: _core_Range | null;
  _targetValue?: "targetValue";
  _rangeValue?: "rangeValue";
}

export interface GoalMetric__Output {
  id: string;
  condition: _core_MetricCondition__Output;
  targetValue?: number;
  rangeValue?: _core_Range__Output | null;
  _targetValue: "targetValue";
  _rangeValue: "rangeValue";
}
