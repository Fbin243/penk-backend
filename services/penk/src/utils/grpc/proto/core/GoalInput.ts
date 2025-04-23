// Original file: ../../proto/core/goal_message.proto

import type { GoalMetricInput as _core_GoalMetricInput, GoalMetricInput__Output as _core_GoalMetricInput__Output } from '../core/GoalMetricInput';
import type { CheckboxInput as _core_CheckboxInput, CheckboxInput__Output as _core_CheckboxInput__Output } from '../core/CheckboxInput';
import type { Long } from '@grpc/proto-loader';

export interface GoalInput {
  'id'?: (string);
  'characterId'?: (string);
  'name'?: (string);
  'description'?: (string);
  'startTime'?: (number | string | Long);
  'endTime'?: (number | string | Long);
  'metrics'?: (_core_GoalMetricInput)[];
  'checkboxes'?: (_core_CheckboxInput)[];
  '_id'?: "id";
  '_description'?: "description";
}

export interface GoalInput__Output {
  'id'?: (string);
  'characterId': (string);
  'name': (string);
  'description'?: (string);
  'startTime': (string);
  'endTime': (string);
  'metrics': (_core_GoalMetricInput__Output)[];
  'checkboxes': (_core_CheckboxInput__Output)[];
  '_id': "id";
  '_description': "description";
}
